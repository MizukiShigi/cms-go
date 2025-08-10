package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/MizukiShigi/cms-go/internal/presentation/helper"
)

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c,omitempty"`
}

func AuthMiddleware(auth0Domain, audience string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			// Authorization ヘッダーからトークンを抽出
			tokenString, err := extractTokenFromHeader(r)
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				myerr := valueobject.NewMyError(valueobject.UnauthorizedCode, err.Error())
				helper.RespondWithError(w, myerr)
				return
			}

			// JWT トークンを解析・検証
			claims, err := validateAuth0Token(tokenString, auth0Domain, audience)
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				myerr := valueobject.NewMyError(valueobject.UnauthorizedCode, err.Error())
				helper.RespondWithError(w, myerr)
				return
			}

			// auth0のカスタムクレームにあるapp_user_idをユーザーIDとして使用する
			claimUserID, ok := claims["app_user_id"].(string)
			if !ok {
				myerr := valueobject.NewMyError(valueobject.UnauthorizedCode, "app_user_id is not a string")
				helper.RespondWithError(w, myerr)
				return
			}

			// 検証されたクレームをコンテキストに追加
			ctx = context.WithValue(r.Context(), domaincontext.UserID, claimUserID)

			// 検証されたクレームをログコンテキストに追加
			ctx = domaincontext.WithValue(ctx, "user_id", claimUserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractTokenFromHeader は Authorization ヘッダーから JWT トークンを抽出
func extractTokenFromHeader(r *http.Request) (string, error) {
	// Authorization ヘッダーを取得
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	// Bearer プレフィックスを確認
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format. 'Bearer <token>' format is required")
	}

	return parts[1], nil
}

// getJWKS はAuth0のJWKSエンドポイントから公開鍵を取得
func getJWKS(auth0Domain string) (*JWKS, error) {
	resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", auth0Domain))
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	return &jwks, nil
}

// getPubKey はJWKから公開鍵を生成
func getPubKey(jwk JWK) (*rsa.PublicKey, error) {
	// base64url-encodedされた値をデコード
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode N: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode E: %w", err)
	}

	// big.Intに変換
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	// RSA公開鍵を作成
	pubKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return pubKey, nil
}

// validateAuth0Token はAuth0のJWTトークンを検証
func validateAuth0Token(tokenString string, auth0Domain string, audience string) (jwt.MapClaims, error) {
	// JWKS取得
	jwks, err := getJWKS(auth0Domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %w", err)
	}

	// トークンを解析（署名検証なし）してヘッダーのkidを取得
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token header: %w", err)
	}

	// kidを取得
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("kid not found in token header")
	}

	// kidに対応するJWKを検索
	var jwk *JWK
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			jwk = &key
			break
		}
	}

	if jwk == nil {
		return nil, fmt.Errorf("JWK not found for kid: %s", kid)
	}

	// 公開鍵を取得
	pubKey, err := getPubKey(*jwk)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key: %w", err)
	}

	// トークンを検証
	validToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// RS256アルゴリズムを厳密に確認
		if token.Method != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("unexpected signing method: expected RS256, got %v", token.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	if !validToken.Valid {
		return nil, errors.New("invalid token")
	}

	// クレームを取得
	claims, ok := validToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to get token claims")
	}

	// issuerを検証
	expectedIssuer := fmt.Sprintf("https://%s/", auth0Domain)
	if iss, ok := claims["iss"].(string); !ok || iss != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer: expected %s, got %v", expectedIssuer, claims["iss"])
	}

	// audienceを検証
	if aud, ok := claims["aud"].(string); ok {
		if aud != audience {
			return nil, fmt.Errorf("invalid audience: expected %s, got %s", audience, aud)
		}
	} else if audSlice, ok := claims["aud"].([]interface{}); ok {
		// audienceが配列の場合
		found := false
		for _, a := range audSlice {
			if audStr, ok := a.(string); ok && audStr == audience {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("invalid audience: expected %s, got %v", audience, claims["aud"])
		}
	} else {
		return nil, fmt.Errorf("invalid audience format: %v", claims["aud"])
	}

	return claims, nil
}

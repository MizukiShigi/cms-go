package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/MizukiShigi/cms-go/infrastructure/auth"
	"github.com/MizukiShigi/cms-go/infrastructure/logger"
	"github.com/MizukiShigi/cms-go/infrastructure/repository"
	"github.com/MizukiShigi/cms-go/internal/presentation/controller"
	"github.com/MizukiShigi/cms-go/internal/presentation/middleware"

	"github.com/MizukiShigi/cms-go/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// ローカル環境用環境変数ファイル読み込み
	loadLocalEnv()

	// ロギング設定
	baseHadler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})
	customHandler := logger.NewHandler(baseHadler)
	slog.SetDefault(slog.New(customHandler))

	// 環境変数の検証とデフォルト値設定
	env := os.Getenv("ENV")
	host := getEnvOrDefault("DB_HOST", "localhost")
	name := getEnvOrDefault("DB_NAME", "cms_dev")
	user := getEnvOrDefault("DB_USER", "postgres")
	password := getEnvOrDefault("DB_PASSWORD", "postgres")
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	port := getEnvOrDefault("PORT", "8080")

	// 必須環境変数の検証
	if env == "" {
		log.Fatal("ENV environment variable is required")
	}
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is required")
	}

	slog.Info("Starting application",
		"db_host", host,
		"db_name", name,
		"db_user", user,
		"port", port,
		"env", os.Getenv("ENV"))

	encodedPassword := url.QueryEscape(password)

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", user, encodedPassword, host, name))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// データベース接続プール設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// DB接続の検証
	if err := db.Ping(); err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}

	// リポジトリ初期化
	transactionManager := repository.NewTransactionManager(db)
	userRepository := repository.NewUserRepository(db)
	postRepository := repository.NewPostRepository(db)
	tagRepository := repository.NewTagRepository(db)

	// サービス初期化
	authService := auth.NewJWTService(jwtSecret)

	// ユースケース初期化
	registerUserUsecase := usecase.NewRegisterUserUsecase(userRepository)
	loginUserUsecase := usecase.NewLoginUserUsecase(userRepository, authService)
	createPostUsecase := usecase.NewCreatePostUsecase(transactionManager, postRepository, tagRepository)
	getPostUsecase := usecase.NewGetPostUsecase(postRepository)
	updatePostUsecase := usecase.NewUpdatePostUsecase(transactionManager, postRepository, tagRepository)
	patchPostUsecase := usecase.NewPatchPostUsecase(postRepository)

	// コントローラー初期化
	authController := controller.NewAuthController(registerUserUsecase, loginUserUsecase)
	postController := controller.NewPostController(createPostUsecase, getPostUsecase, updatePostUsecase, patchPostUsecase)

	// ルーティング設定
	r := mux.NewRouter()

	// 全てのリクエストにミドルウェア設定
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.TimeoutMiddleware)

	// バージョニング
	v1Router := r.PathPrefix("/cms/v1").Subrouter()

	// 認証不要エンドポイント
	publicV1Router := v1Router.PathPrefix("/").Subrouter()

	// 認証
	authRouter := publicV1Router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", authController.Register).Methods("POST")
	authRouter.HandleFunc("/login", authController.Login).Methods("POST")

	// 認証必須エンドポイント
	protectedV1Router := v1Router.PathPrefix("/").Subrouter()
	protectedV1Router.Use(middleware.AuthMiddleware(jwtSecret))

	// 投稿
	postRouter := protectedV1Router.PathPrefix("/posts").Subrouter()
	postRouter.HandleFunc("/", postController.CreatePost).Methods("POST")
	postRouter.HandleFunc("/{id}", postController.GetPost).Methods("GET")
	postRouter.HandleFunc("/{id}", postController.UpdatePost).Methods("PUT")
	postRouter.HandleFunc("/{id}", postController.PatchPost).Methods("PATCH")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// シグナルハンドリング
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// サーバー起動
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("server is running on port %s\n", port)

	// シグナルを受け取り、コンテキストをキャンセルする
	<-ctx.Done()
	log.Println("server is shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %s\n", err)
	}

	log.Println("server exited properly")
}

// getEnvOrDefault は環境変数を取得し、設定されていない場合はデフォルト値を返す
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadLocalEnv() {
	env := os.Getenv("ENV")
	if env == "local" || env == "" {
		// ローカル環境のみ .env ファイルを読み込む
		if err := godotenv.Load(".env.development"); err != nil {
			log.Fatalf(".env.development ファイルが見つかりません: %v", err)
		}
	}
}

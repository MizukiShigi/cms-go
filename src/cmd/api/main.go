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
	// 開発環境用環境変数ファイル読み込み
	loadDevelopEnv()

	// ロギング設定
	baseHadler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})
	customHandler := logger.NewHandler(baseHadler)
	slog.SetDefault(slog.New(customHandler))

	// DBセットアップ
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	if host == "" || name == "" || user == "" || password == "" {
		log.Fatal("Database connection environment variables not set")
	}

	encodedPassword := url.QueryEscape(password)

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", user, encodedPassword, host, name))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
	authService := auth.NewJWTService(os.Getenv("JWT_SECRET_KEY"))

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
	protectedV1Router.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")))

	// 投稿
	postRouter := protectedV1Router.PathPrefix("/posts").Subrouter()
	postRouter.HandleFunc("/", postController.CreatePost).Methods("POST")
	postRouter.HandleFunc("/{id}", postController.GetPost).Methods("GET")
	postRouter.HandleFunc("/{id}", postController.UpdatePost).Methods("PUT")
	postRouter.HandleFunc("/{id}", postController.PatchPost).Methods("PATCH")

	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
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
	log.Printf("server is running on port %s\n", os.Getenv("PORT"))

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

func loadDevelopEnv() {
	env := os.Getenv("GO_ENV")
	if env == "" || env == "development" {
		// 開発環境のみ .env ファイルを読み込む
		if err := godotenv.Load(".env.development"); err != nil {
			log.Printf(".env.development ファイルが見つかりません（本番環境では正常）: %v", err)
			// log.Fatalf → log.Printf に変更してアプリを継続
		}
	}
}

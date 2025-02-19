package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/config"
	"github.com/223n-tech/SuiminNisshi-Go/internal/handler"
	"github.com/223n-tech/SuiminNisshi-Go/internal/middleware"
	"github.com/223n-tech/SuiminNisshi-Go/internal/repository"
	"github.com/223n-tech/SuiminNisshi-Go/internal/service"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// ロガーの初期化
	logger := log.New(os.Stdout, "[SuiminNisshi] ", log.LstdFlags|log.Lshortfile)

	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// データベース接続の初期化
	db, err := repository.NewDB(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// テンプレートマネージャーの初期化
	tm := handler.NewTemplateManager("web/template", nil)
	if err := tm.LoadTemplates(); err != nil {
		logger.Fatalf("Failed to load templates: %v", err)
	}

	// リポジトリの初期化
	repo := repository.NewRepository(db)

	// サービスの初期化
	svc := service.NewService(repo)

	// ルーターの設定
	r := chi.NewRouter()

	// ミドルウェアの設定
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// カスタムミドルウェアの設定
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.SecurityHeaders)

	// 静的ファイルの提供
	fileServer := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// ルートの登録
	router := handler.NewRouter(r)

	// サービスをハンドラーに渡す
	authHandler := handler.NewAuthHandler(tm)
	authHandler.RegisterRoutes(router)

	// ダッシュボードハンドラーの初期化と登録
	dashboardHandler := handler.NewDashboardHandler(tm, svc)
	dashboardHandler.RegisterRoutes(router)

	// プロフィールハンドラーの初期化と登録
	profileHandler := handler.NewProfileHandler(tm)
	profileHandler.RegisterRoutes(router)

	// 設定ハンドラーの初期化と登録
	settingsHandler := handler.NewSettingsHandler(tm)
	settingsHandler.RegisterRoutes(router)

	// エラーハンドラーの初期化
	errorHandler, err := handler.NewErrorHandler()
	if err != nil {
		logger.Fatalf("Failed to initialize error handler: %v", err)
	}
	// 404ハンドラーの設定
	r.NotFound(errorHandler.Handle404)
	// // ミドルウェアでパニックをキャッチして500エラーを表示
	// r.Use(func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		defer func() {
	// 			if err := recover(); err != nil {
	// 				errorHandler.Handle500(w, r, fmt.Errorf("%v", err))
	// 			}
	// 		}()
	// 		next.ServeHTTP(w, r)
	// 	})
	// })

	// サーバーの設定
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// サーバーの起動（ゴルーチンで実行）
	go func() {
		logger.Printf("Server is starting on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// グレースフルシャットダウンの設定
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Server is shutting down...")

	// シャットダウンのコンテキスト
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// サーバーのシャットダウン
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server stopped gracefully")
}

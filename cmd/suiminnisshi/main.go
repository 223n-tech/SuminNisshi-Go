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
	logger.Printf("[Initialize] SuiminNisshi Startup...")

	// 設定の読み込み
	logger.Printf("[Initialize] Loading config...")
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("[NG] Failed to load config: %v", err)
	}

	// データベース接続の初期化
	logger.Printf("[Initialize] Connecting to database...")
	db, err := repository.NewDB(cfg.Database)
	if err != nil {
		logger.Fatalf("[NG] Failed to connect to database: %v", err)
	}
	defer db.Close()

	// テンプレートマネージャーの初期化
	logger.Printf("[Initialize] Loading templates...")
	tm := handler.NewTemplateManager("web/views", nil)
	if err := tm.LoadTemplates(); err != nil {
		logger.Fatalf("[NG] Failed to load templates: %v", err)
	}

	// リポジトリの初期化
	logger.Printf("[Initialize] Initializing repository...")
	repo := repository.NewRepository(db)

	// サービスの初期化
	logger.Printf("[Initialize] Initializing service...")
	svc := service.NewService(repo)

	// ルーターの設定
	logger.Printf("[Initialize] Setting up router...")
	r := chi.NewRouter()

	// ミドルウェアの設定
	logger.Printf("[Initialize] Setting up middleware...")
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
	logger.Printf("[Initialize] Setting up custom middleware...")
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.SecurityHeaders)

	// 静的ファイルの提供
	logger.Printf("[Initialize] Setting up static file server...")
	fileServer := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// ルートの登録
	logger.Printf("[Initialize] Registering routes...")
	router := handler.NewRouter(r)

	// サービスをハンドラーに渡す
	logger.Printf("[Initialize] Passing service to handlers...")
	authHandler := handler.NewAuthHandler(tm)
	authHandler.RegisterRoutes(router)

	// ダッシュボードハンドラーの初期化と登録
	logger.Printf("[Initialize] Registering dashboard routes...")
	dashboardHandler := handler.NewDashboardHandler(tm, svc)
	dashboardHandler.RegisterRoutes(router)

	// サーバーの設定
	logger.Printf("[Initialize] Setting up server...")
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// サーバーの起動（ゴルーチンで実行）
	go func() {
		logger.Printf("[START] Server is starting on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("[NG] Failed to start server: %v", err)
		}
	}()

	// グレースフルシャットダウンの設定
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("[STOP] Server is shutting down...")

	// シャットダウンのコンテキスト
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// サーバーのシャットダウン
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("[STOP] Server forced to shutdown: %v", err)
	}

	logger.Println("[STOP] Server stopped gracefully")
}

// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/error.go
// errorは、エラーページのハンドラーを提供します。

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
)

// エラーページのハンドラー
type ErrorHandler struct {
	templates *TemplateManager
	service   *service.Service
	logger    *log.Logger
}

// ErrorHandlerを作成
func NewErrorHandler(templates *TemplateManager, svc *service.Service, logger *log.Logger) *ErrorHandler {
	return &ErrorHandler{
		templates: templates,
		service:   svc,
		logger:    logger,
	}
}

// エラーページを表示
func (h *ErrorHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request, status int, err error) {
	// エラーのログ記録
	if err != nil {
		h.logger.Printf("Error occurred: %v\nStack trace:\n%s", err, debug.Stack())
	}

	// ステータスコードに基づいてテンプレート名を決定
	templateName := h.getTemplateNameForStatus(status)

	// エラーメッセージの準備
	message := h.getErrorMessage(status, err)

	// レスポンスヘッダーの設定
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	// テンプレートデータの準備
	data := &TemplateData{
		Title: fmt.Sprintf("%d - %s", status, http.StatusText(status)),
		Data: map[string]interface{}{
			"Status":      status,
			"StatusText": http.StatusText(status),
			"Message":    message,
		},
	}

	// テンプレートのレンダリング
	if err := h.templates.Render(w, templateName, data); err != nil {
		h.logger.Printf("Error rendering error template: %v", err)
		http.Error(w, http.StatusText(status), status)
	}
}

// 404エラーページを表示
func (h *ErrorHandler) Handle404(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r, http.StatusNotFound, nil)
}

// 500エラーページを表示
func (h *ErrorHandler) Handle500(w http.ResponseWriter, r *http.Request, err error) {
	h.ServeHTTP(w, r, http.StatusInternalServerError, err)
}

// 403エラーページを表示
func (h *ErrorHandler) Handle403(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r, http.StatusForbidden, nil)
}

// 405エラーページを表示
func (h *ErrorHandler) Handle405(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r, http.StatusMethodNotAllowed, nil)
}

// パニック時のリカバリー処理
func (h *ErrorHandler) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				h.logger.Printf("Panic recovered: %v\nStack trace:\n%s", err, debug.Stack())
				h.Handle500(w, r, fmt.Errorf("panic: %v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ステータスコードに対応するテンプレート名を取得
func (h *ErrorHandler) getTemplateNameForStatus(status int) string {
	switch status {
	case http.StatusNotFound:
		return "404.html"
	case http.StatusForbidden:
		return "403.html"
	case http.StatusMethodNotAllowed:
		return "405.html"
	default:
		return "500.html"
	}
}

// エラーメッセージを取得
func (h *ErrorHandler) getErrorMessage(status int, err error) string {
	switch status {
	case http.StatusNotFound:
		return "お探しのページが見つかりませんでした。"
	case http.StatusForbidden:
		return "このページにアクセスする権限がありません。"
	case http.StatusMethodNotAllowed:
		return "許可されていないメソッドです。"
	case http.StatusInternalServerError:
		if err != nil {
			return fmt.Sprintf("サーバーエラーが発生しました：%v", err)
		}
		return "サーバーエラーが発生しました。"
	default:
		return http.StatusText(status)
	}
}

// エラーをログに記録
func (h *ErrorHandler) LogError(r *http.Request, err error) {
	h.logger.Printf("Error handling request for %s: %v", r.URL.Path, err)
}

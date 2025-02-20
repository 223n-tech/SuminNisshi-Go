package handler

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// ErrorHandler エラーページのハンドラー
type ErrorHandler struct {
	templates map[int]*template.Template
}

// NewErrorHandler エラーハンドラーを作成
func NewErrorHandler() (*ErrorHandler, error) {
	handler := &ErrorHandler{
		templates: make(map[int]*template.Template),
	}

	// エラーテンプレートの読み込み
	errorPages := map[int]string{
		http.StatusNotFound:            "404.html",
		http.StatusInternalServerError: "500.html",
		http.StatusForbidden:          "403.html",
	}

	for status, filename := range errorPages {
		template, err := template.ParseFiles(filepath.Join("web", "template", "errors", filename))
		if err != nil {
			return nil, err
		}
		handler.templates[status] = template
	}

	return handler, nil
}

// ServeHTTP エラーページを表示
func (h *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, status int) {
	template, exists := h.templates[status]
	if !exists {
		log.Printf("No template for status %d, falling back to 500", status)
		template = h.templates[http.StatusInternalServerError]
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
	if err := template.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

// Handle404 404エラーページを表示
func (h *ErrorHandler) Handle404(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r, http.StatusNotFound)
}

// Handle500 500エラーページを表示
func (h *ErrorHandler) Handle500(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal Server Error: %v", err)
	h.ServeHTTP(w, r, http.StatusInternalServerError)
}

// Handle403 403エラーページを表示
func (h *ErrorHandler) Handle403(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r, http.StatusForbidden)
}

package handler

import (
	"log"
	"net/http"
)

// ErrorHandler エラーページのハンドラー
type ErrorHandler struct {
	templates *TemplateManager
}

// NewErrorHandlerWithTemplates エラーハンドラーを作成
func NewErrorHandlerWithTemplates(templates *TemplateManager) *ErrorHandler {
	return &ErrorHandler{
		templates: templates,
	}
}

// ServeHTTP エラーページを表示
func (h *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, status int) {
	data := &TemplateData{
		Title: http.StatusText(status),
	}

	var templateName string
	switch status {
	case http.StatusNotFound:
		templateName = "404.html"
	case http.StatusForbidden:
		templateName = "403.html"
	case http.StatusMethodNotAllowed:
		templateName = "405.html"
	default:
		templateName = "500.html"
	}

	err := h.templates.Render(w, templateName, data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", templateName, err)
		http.Error(w, http.StatusText(status), status)
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

// Handle405 405エラーページを表示
func (h *ErrorHandler) Handle405(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r, http.StatusMethodNotAllowed)
}

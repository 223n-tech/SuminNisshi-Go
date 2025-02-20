package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// PrivacyHandler プライバシーポリシー画面のハンドラー
type PrivacyHandler struct {
	templates *TemplateManager
}

// NewPrivacyHandler プライバシーポリシーハンドラーを作成
func NewPrivacyHandler(templates *TemplateManager) *PrivacyHandler {
	return &PrivacyHandler{
		templates: templates,
	}
}

// RegisterRoutes ルートの登録
func (h *PrivacyHandler) RegisterRoutes(r chi.Router) {
	r.Get("/privacy", h.Privacy)
}

// Privacy プライバシーポリシーページの表示
func (h *PrivacyHandler) Privacy(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "プライバシーポリシー",
	}

	err := h.templates.Render(w, "privacy.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

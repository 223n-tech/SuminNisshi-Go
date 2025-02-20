// プライバシーポリシーページのハンドラー
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

/*
	PrivacyHandler プライバシーポリシーページのハンドラ
*/
type PrivacyHandler struct {
	templates *TemplateManager
}

/*
	NewPrivacyHandler は PrivacyHandler を作成します。
*/
func NewPrivacyHandler(templates *TemplateManager) *PrivacyHandler {
	return &PrivacyHandler{
		templates: templates,
	}
}

/*
	RegisterRoutes ルーティングを登録
*/
func (h *PrivacyHandler) RegisterRoutes(r chi.Router) {
	r.Get("/privacy", h.Privacy)
}

/*
	Privacy プライバシーポリシーページを表示
*/
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

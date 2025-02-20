package handler

import (
    "net/http"

    "github.com/go-chi/chi/v5"
)

// TermsHandler 利用規約画面のハンドラー
type TermsHandler struct {
    templates *TemplateManager
}

// NewTermsHandler 利用規約ハンドラーを作成
func NewTermsHandler(templates *TemplateManager) *TermsHandler {
    return &TermsHandler{
        templates: templates,
    }
}

// RegisterRoutes ルートの登録
func (h *TermsHandler) RegisterRoutes(r chi.Router) {
    r.Get("/terms", h.Terms)
}

// Terms 利用規約ページの表示
func (h *TermsHandler) Terms(w http.ResponseWriter, r *http.Request) {
    data := &TemplateData{
        Title: "利用規約",
    }

    err := h.templates.Render(w, "terms.html", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

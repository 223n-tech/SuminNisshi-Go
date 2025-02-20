// 利用規約画面のハンドラーを定義
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

/*
    TermsHandler は利用規約画面のハンドラです。
*/
type TermsHandler struct {
    templates *TemplateManager
}

/*
    NewTermsHandler は TermsHandler を作成します。
*/
func NewTermsHandler(templates *TemplateManager) *TermsHandler {
    return &TermsHandler{
        templates: templates,
    }
}

/*
    RegisterRoutes はルーティングを登録します。
*/
func (h *TermsHandler) RegisterRoutes(r chi.Router) {
    r.Get("/terms", h.Terms)
}

/*
    Terms 利用規約画面を表示
*/
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

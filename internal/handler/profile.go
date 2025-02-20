// プロフィール画面のハンドラー
package handler

import (
	"net/http"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

/*
	ProfileHandler プロフィール関連のハンドラ
*/
type ProfileHandler struct {
	templates *TemplateManager
}

/*
	NewProfileHandler は ProfileHandler を作成します。
*/
func NewProfileHandler(templates *TemplateManager) *ProfileHandler {
	return &ProfileHandler{
		templates: templates,
	}
}

/*
	RegisterRoutes ルーティングを登録
*/
func (h *ProfileHandler) RegisterRoutes(r *RouterWrapper) {
	r.Get("/profile", h.Profile)
	r.Post("/profile/update", h.UpdateProfile)
}

/*
	Profile プロフィール画面を表示
*/
func (h *ProfileHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// TODO: 実際のユーザー情報を取得
	data := &TemplateData{
		Title:      "プロフィール",
		ActiveMenu: "profile",
		User: &models.User{
			Name:  "テストユーザー",
			Email: "test@example.com",
			Timezone: "Asia/Tokyo",
		},
	}

	err := h.templates.Render(w, "profile.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	UpdateProfile プロフィールの更新
*/
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: プロフィール更新の実装
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/privacy.go
// privacyは、プライバシーポリシーページのハンドラーを提供します。

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
	"github.com/go-chi/chi/v5"
)

// プライバシーポリシーページのハンドラー
type PrivacyHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// PrivacyHandlerを作成
func NewPrivacyHandler(templates *TemplateManager, svc *service.Service) *PrivacyHandler {
	return &PrivacyHandler{
		templates: templates,
		service:   svc,
	}
}

// ルーティングを登録
func (h *PrivacyHandler) RegisterRoutes(r chi.Router) {
	r.Get("/privacy", h.Privacy)
}

// プライバシーポリシーページを表示
func (h *PrivacyHandler) Privacy(w http.ResponseWriter, _ *http.Request) {
	// ユーザー情報の取得（オプション）

	// ダミーユーザー情報
	user := &models.User{
		ID:                1,
		Email:             "test@example.com",
		DisplayName:       "テストユーザー",
		LastLoginDatetime: sql.NullTime{ Time: time.Now(), Valid: true },
		Created:           time.Now(),
		Modified:          time.Now(),
		Deleted:           sql.NullTime{},
	}

	// TODO: UserKeyを使ってユーザー情報を取得
	// if userCtx := r.Context().Value(UserKey); userCtx != nil {
	// 	userID = GetUserIDFromContext(userCtx)
	// 	user, err := h.service.User().GetUserByID(r.Context(), userID)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	// プライバシーポリシーのメタデータ
	metadata := map[string]interface{}{
		"lastUpdated": time.Date(2025, 2, 22, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
		"version":     "0.1.0",
	}

	data := &TemplateData{
		Title:      "プライバシーポリシー",
		ActiveMenu: "privacy",
		User:       user,
		Data: map[string]interface{}{
			"Metadata": metadata,
		},
		Meta: map[string]interface{}{
			"description": "SuiminNisshi-Goのプライバシーポリシー。個人情報の取り扱いについて説明します。",
			"robots":      "noindex, nofollow", // 検索エンジンにインデックスさせない
		},
	}

	err := h.templates.Render(w, "privacy.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

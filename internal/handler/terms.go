// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/terms.go
// termsは、利用規約画面のハンドラーを提供します。

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
	"github.com/go-chi/chi/v5"
)

// 利用規約画面のハンドラー
type TermsHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// TermsHandlerを作成
func NewTermsHandler(templates *TemplateManager, svc *service.Service) *TermsHandler {
	return &TermsHandler{
		templates: templates,
		service:   svc,
	}
}

// ルーティングを登録
func (h *TermsHandler) RegisterRoutes(r chi.Router) {
	r.Get("/terms", h.Terms)
	r.Get("/terms/history", h.TermsHistory)
}

// 利用規約画面を表示
func (h *TermsHandler) Terms(w http.ResponseWriter, _ *http.Request) {
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
	// 	user = userCtx
	// }

	// 利用規約のメタデータ
	metadata := map[string]interface{}{
		"lastUpdated": time.Date(2025, 2, 22, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
		"version":     "0.1.0",
		"effectiveDate": time.Date(2025, 2, 22, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
	}

	// 利用規約の更新履歴
	history := []map[string]interface{}{
		{
			"version":     "0.1.0",
			"date":        time.Date(2025, 2, 22, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
			"description": "初版リリース",
		},
	}

	
	data := &TemplateData{
		Title:		"利用規約",
		ActiveMenu:	"terms",
		User:		user,
		Data:		map[string]interface{}{
						"Metadata": metadata,
						"History":  history,
					},
		Meta: map[string]interface{}{
			"description":	"SuiminNisshi-Goの利用規約。サービスの利用条件について説明します。",
			"robots":		"noindex, nofollow", // 検索エンジンにインデックスさせない
		},
	}

	err := h.templates.Render(w, "terms.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 利用規約の更新履歴を表示
func (h *TermsHandler) TermsHistory(w http.ResponseWriter, _ *http.Request) {
	// 更新履歴データ
	history := []map[string]interface{}{
		{
			"version":     "0.1.0",
			"date":        time.Date(2025, 2, 22, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
			"description": "初版リリース",
			"changes": []string{
				"サービス全般に関する規約を制定",
				"個人情報の取り扱いに関する規定を追加",
				"利用制限に関する規定を追加",
			},
		},
	}

	data := &TemplateData{
		Title:      "利用規約の更新履歴",
		ActiveMenu: "terms",
		Data: map[string]interface{}{
			"History": history,
		},
		Meta: map[string]interface{}{
			"description": "SuiminNisshi-Goの利用規約の更新履歴",
			"robots":      "noindex, nofollow",
		},
	}

	err := h.templates.Render(w, "terms-history.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

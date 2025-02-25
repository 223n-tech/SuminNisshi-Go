// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/account_deletion.go
// account_deletionは、アカウント削除画面のハンドラーを提供します。

import (
	"context"
	"net/http"

	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
	"github.com/go-chi/chi/v5"
)

// アカウント削除画面のハンドラー
type AccountDeletionHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// アカウント削除画面のハンドラーを作成
func NewAccountDeletionHandler(templates *TemplateManager, svc *service.Service) *AccountDeletionHandler {
	return &AccountDeletionHandler{
		templates: templates,
		service:   svc,
	}
}

// ルーティングを登録
func (h *AccountDeletionHandler) RegisterRoutes(r chi.Router) {
	r.Get("/settings/account/delete", h.ShowDeleteConfirmation)
	r.Post("/settings/account/delete", h.DeleteAccount)
}

// コンテキストからユーザーIDを取得
func GetUserIDFromContext(ctx context.Context) int64 {
	return ctx.Value(UserKey).(int64)
}

// アカウント削除確認画面を表示
func (h *AccountDeletionHandler) ShowDeleteConfirmation(w http.ResponseWriter, r *http.Request) {
	// コンテキストからユーザーIDを取得
	userID := GetUserIDFromContext(r.Context())

	// ユーザー情報の取得
	user, err := h.service.User().GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		data := &TemplateData{
			Title:      "アカウント削除",
			ActiveMenu: "settings",
			Flash: &Flash{
				Type:    "danger",
				Message: "ユーザーが見つかりません",
			},
		}
		h.templates.Render(w, "account-deletion.html", data)
		return
	}

	// TODO: ユーザー削除処理の実装

	data := &TemplateData{
		Title:      "アカウント削除",
		ActiveMenu: "settings",
	}

	err = h.templates.Render(w, "account-deletion.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// アカウント削除の処理
func (h *AccountDeletionHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// コンテキストからユーザーIDを取得
	// userID = GetUserIDFromContext(r.Context())

	// パスワードの検証
	password := r.FormValue("password")
	if password == "" {
		data := &TemplateData{
			Title:      "アカウント削除",
			ActiveMenu: "settings",
			Flash: &Flash{
				Type:    "danger",
				Message: "パスワードを入力してください",
			},
		}
		h.templates.Render(w, "account-deletion.html", data)
		return
	}

	// パスワードの検証
	// TODO: 実際のユーザーIDを使用
	/*
	if err := h.service.User().ValidatePassword(r.Context(), userID, password); err != nil {
		data := &TemplateData{
			Title:      "アカウント削除",
			ActiveMenu: "settings",
			Flash: &Flash{
				Type:    "danger",
				Message: "パスワードが正しくありません",
			},
		}
		h.templates.Render(w, "account-deletion.html", data)
		return
	}
	*/

	// 確認チェックボックスの検証
	if r.FormValue("confirm") != "on" {
		data := &TemplateData{
			Title:      "アカウント削除",
			ActiveMenu: "settings",
			Flash: &Flash{
				Type:    "danger",
				Message: "削除の確認が必要です",
			},
		}
		h.templates.Render(w, "account-deletion.html", data)
		return
	}

	// フィードバックの保存（任意）
	// feedback := r.FormValue("feedback")
	// TODO: フィードバックの保存処理
	// if feedback != "" {

	// }

	// アカウント削除の実行
	// TODO: 実際のユーザーIDを使用
	/*
	if err := h.service.User().DeleteAccount(r.Context(), userID); err != nil {
		data := &TemplateData{
			Title:      "アカウント削除",
			ActiveMenu: "settings",
			Flash: &Flash{
				Type:    "danger",
				Message: "アカウントの削除に失敗しました",
			},
		}
		h.templates.Render(w, "account-deletion.html", data)
		return
	}
	*/

	// セッションの破棄
	// TODO: セッション管理の実装
	// ClearSession(r.Context())

	// 完了ページにリダイレクト
	http.Redirect(w, r, "/account-deleted", http.StatusSeeOther)
}

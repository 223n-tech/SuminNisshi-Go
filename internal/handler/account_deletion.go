package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// AccountDeletionHandler アカウント削除画面のハンドラー
type AccountDeletionHandler struct {
	templates *TemplateManager
}

// NewAccountDeletionHandler アカウント削除ハンドラーを作成
func NewAccountDeletionHandler(templates *TemplateManager) *AccountDeletionHandler {
	return &AccountDeletionHandler{
		templates: templates,
	}
}

// RegisterRoutes ルートの登録
func (h *AccountDeletionHandler) RegisterRoutes(r chi.Router) {
	r.Get("/settings/account/delete", h.ShowDeleteConfirmation)
	r.Post("/settings/account/delete", h.DeleteAccount)
}

// ShowDeleteConfirmation アカウント削除確認ページの表示
func (h *AccountDeletionHandler) ShowDeleteConfirmation(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "アカウント削除",
		// CSRFToken: GetCSRFToken(r), // CSRFトークン生成関数は別途実装が必要
	}

	err := h.templates.Render(w, "account-deletion.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteAccount アカウント削除の実行
func (h *AccountDeletionHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	// POSTデータの解析
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// パスワードの検証
	// password := r.FormValue("password")
	// if !ValidatePassword(r.Context(), password) { // パスワード検証関数は別途実装が必要
	// 	http.Error(w, "パスワードが正しくありません", http.StatusUnauthorized)
	// 	return
	// }

	// 確認チェックボックスの検証
	if r.FormValue("confirm") != "on" {
		http.Error(w, "削除の確認が必要です", http.StatusBadRequest)
		return
	}

	// フィードバックの保存（任意）
	feedback := r.FormValue("feedback")
	if feedback != "" {
		// SaveFeedback(r.Context(), feedback) // フィードバック保存関数は別途実装が必要
	}

	// アカウント削除の実行
	// if err := DeleteUserAccount(r.Context()); err != nil { // アカウント削除関数は別途実装が必要
	// 	http.Error(w, "アカウントの削除に失敗しました", http.StatusInternalServerError)
	// 	return
	// }

	// セッションの破棄
	// ClearSession(r.Context()) // セッション破棄関数は別途実装が必要

	// 完了ページにリダイレクト
	http.Redirect(w, r, "/account-deleted", http.StatusSeeOther)
}

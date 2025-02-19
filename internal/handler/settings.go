package handler

import (
	"net/http"
)

// SettingsHandler 設定画面のハンドラー
type SettingsHandler struct {
	templates *TemplateManager
}

// NewSettingsHandler 設定ハンドラーを作成
func NewSettingsHandler(templates *TemplateManager) *SettingsHandler {
	return &SettingsHandler{
		templates: templates,
	}
}

// RegisterRoutes ルートの登録
func (h *SettingsHandler) RegisterRoutes(r *RouterWrapper) {
	r.Get("/settings", h.Settings)
	r.Post("/settings/profile", h.UpdateProfile)
	r.Post("/settings/password", h.UpdatePassword)
	r.Post("/settings/notifications", h.UpdateNotifications)
	r.Post("/settings/export/csv", h.ExportCSV)
	r.Post("/settings/export/json", h.ExportJSON)
	r.Post("/settings/account/delete", h.DeleteAccount)
}

// Settings 設定ページの表示
func (h *SettingsHandler) Settings(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:      "設定",
		ActiveMenu: "settings",
		User: &User{
			Name:  "テストユーザー",
			Email: "test@example.com",
		},
	}

	err := h.templates.Render(w, "settings.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateProfile プロフィール更新
func (h *SettingsHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: プロフィール更新の実装
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// UpdatePassword パスワード更新
func (h *SettingsHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	// TODO: パスワード更新の実装
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// UpdateNotifications 通知設定更新
func (h *SettingsHandler) UpdateNotifications(w http.ResponseWriter, r *http.Request) {
	// TODO: 通知設定更新の実装
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// ExportCSV CSVエクスポート
func (h *SettingsHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	// TODO: CSVエクスポートの実装
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// ExportJSON JSONエクスポート
func (h *SettingsHandler) ExportJSON(w http.ResponseWriter, r *http.Request) {
	// TODO: JSONエクスポートの実装
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// DeleteAccount アカウント削除
func (h *SettingsHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	// TODO: アカウント削除の実装
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

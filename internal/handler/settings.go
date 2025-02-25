// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/settings.go
// settingsは、設定画面のハンドラーを提供します。

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
	"github.com/go-chi/chi/v5"
)

// 設定ページのハンドラー
type SettingsHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// SettingsHandlerを作成
func NewSettingsHandler(templates *TemplateManager, svc *service.Service) *SettingsHandler {
	return &SettingsHandler{
		templates: templates,
		service:   svc,
	}
}

// ルーティングを登録
func (h *SettingsHandler) RegisterRoutes(r chi.Router) {
	r.Get("/settings", h.Settings)
	r.Post("/settings/profile", h.UpdateProfile)
	r.Post("/settings/password", h.UpdatePassword)
	r.Post("/settings/notifications", h.UpdateNotifications)
	r.Get("/settings/export/csv", h.ExportCSV)
	r.Get("/settings/export/json", h.ExportJSON)
	r.Get("/settings/account/delete", h.ShowDeleteAccountPage)
	r.Post("/settings/account/delete", h.DeleteAccount)
}

// 設定画面を表示
func (h *SettingsHandler) Settings(w http.ResponseWriter, r *http.Request) {
	// ユーザー情報の取得
	userCtx := r.Context().Value(UserKey)
	if userCtx == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// ユーザー情報の取得
	user, err := h.service.User().GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "ユーザー情報の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// 睡眠設定の取得
	pref, err := h.service.User().GetSleepPreference(r.Context(), userID)
	if err != nil {
		http.Error(w, "睡眠設定の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	data := &TemplateData{
		Title:      "設定",
		ActiveMenu: "settings",
		User:       user,
		Data: map[string]interface{}{
			"Preferences": pref,
		},
	}

	// フラッシュメッセージの設定
	if msg := r.URL.Query().Get("message"); msg != "" {
		data.Flash = &Flash{
			Type:    r.URL.Query().Get("type"),
			Message: msg,
		}
	}

	err = h.templates.Render(w, "settings.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// プロフィールの更新
func (h *SettingsHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// プロフィール情報の更新
	user := &models.User{
		ID:          userID,
		DisplayName: r.FormValue("display_name"),
		Email:       r.FormValue("email"),
		TimeZone:    r.FormValue("timezone"),
	}

	err := h.service.User().UpdateProfile(r.Context(), user)
	if err != nil {
		http.Redirect(w, r, "/settings?message=プロフィールの更新に失敗しました&type=danger", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/settings?message=プロフィールを更新しました&type=success", http.StatusSeeOther)
}

// パスワードの更新
func (h *SettingsHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")

	err := h.service.User().UpdatePassword(r.Context(), userID, currentPassword, newPassword)
	if err != nil {
		http.Redirect(w, r, "/settings?message=パスワードの更新に失敗しました&type=danger", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/settings?message=パスワードを更新しました&type=success", http.StatusSeeOther)
}

// 通知設定の更新
func (h *SettingsHandler) UpdateNotifications(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	pref := &models.UserSleepPreference{
		UserID:            userID,
		IsReminderEnabled: r.FormValue("reminder_enabled") == "on",
	}

	err := h.service.User().UpdateSleepPreference(r.Context(), pref)
	if err != nil {
		http.Redirect(w, r, "/settings?message=通知設定の更新に失敗しました&type=danger", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/settings?message=通知設定を更新しました&type=success", http.StatusSeeOther)
}

// CSVエクスポート
func (h *SettingsHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
    // TODO: 実際のユーザーIDを使用
    var userID int64 = 1 // 開発用

    // CSVファイル名の設定
    filename := fmt.Sprintf("sleep-records-%s.csv", time.Now().Format("2006-01-02"))
    w.Header().Set("Content-Type", "text/csv")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

    writer := csv.NewWriter(w)
    defer writer.Flush()

    // ヘッダーの書き込み
    headers := []string{"日付", "時間枠", "状態", "種別", "メモ"}
    if err := writer.Write(headers); err != nil {
        http.Error(w, "CSVの書き込みに失敗しました", http.StatusInternalServerError)
        return
    }

    // データの取得と書き込み
    records, err := h.service.Record().GetAllRecords(r.Context(), userID)
    if err != nil {
        http.Error(w, "データの取得に失敗しました", http.StatusInternalServerError)
        return
    }

    for _, record := range records {
        // 実際のSleepRecord構造体のフィールドに基づいて行を作成
        row := []string{
            record.RecordDate.Format("2006-01-02"),
            record.TimeSlot.Format("15:04"),
            fmt.Sprintf("%d", record.SleepStateID), // 本来は状態名を取得すべき
            record.RecordType,
            record.Note.String,
        }
        if err := writer.Write(row); err != nil {
            http.Error(w, "CSVの書き込みに失敗しました", http.StatusInternalServerError)
            return
        }
    }
}

// JSONエクスポート
func (h *SettingsHandler) ExportJSON(w http.ResponseWriter, r *http.Request) {
	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// JSONファイル名の設定
	filename := fmt.Sprintf("sleep-records-%s.json", time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// データの取得
	records, err := h.service.Record().GetAllRecords(r.Context(), userID)
	if err != nil {
		http.Error(w, "データの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// JSONエンコード
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(records); err != nil {
		http.Error(w, "JSONの書き込みに失敗しました", http.StatusInternalServerError)
		return
	}
}

// アカウント削除の確認画面を表示
func (h *SettingsHandler) ShowDeleteAccountPage(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:      "アカウント削除",
		ActiveMenu: "settings",
	}

	err := h.templates.Render(w, "delete-account.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// アカウントの削除
func (h *SettingsHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// password := r.FormValue("password")
	confirm := r.FormValue("confirm") == "on"

	if !confirm {
		http.Redirect(w, r, "/settings/account/delete?message=削除の確認が必要です&type=danger", http.StatusSeeOther)
		return
	}

	err := h.service.User().DeleteAccount(r.Context(), userID)
	if err != nil {
		http.Redirect(w, r, "/settings/account/delete?message=アカウントの削除に失敗しました&type=danger", http.StatusSeeOther)
		return
	}

	// セッションのクリア
	// TODO: セッション管理の実装
	// ClearSession(r)

	http.Redirect(w, r, "/account-deleted", http.StatusSeeOther)
}

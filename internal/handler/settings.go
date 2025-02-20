// 設定画面のハンドラー
package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
	"github.com/go-chi/chi/v5"
)

/*
	SettingsHandler 設定ページのハンドラ
*/
type SettingsHandler struct {
	templates *TemplateManager
}

/*
	User ユーザー情報の構造体
*/
type User struct {
	Name     string
	Email    string
	Timezone string
}

/*
	NotificationSettings 通知設定の構造体
*/
type NotificationSettings struct {
	EmailEnabled    bool
	BedtimeReminder bool
	ReminderTime    string
	WeeklyReport    bool
}

/*
	NewSettingsHandler は SettingsHandler を作成します。
*/
func NewSettingsHandler(templates *TemplateManager) *SettingsHandler {
	return &SettingsHandler{
		templates: templates,
	}
}

/*
	RegisterRoutes ルーティングを登録
*/
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

/*
	Settings 設定画面を表示
*/
func (h *SettingsHandler) Settings(w http.ResponseWriter, r *http.Request) {
	// テストデータの作成
	user := &models.User{
		Name:     "テストユーザー",
		Email:    "test@example.com",
		Timezone: "Asia/Tokyo",
	}

	notifications := &models.NotificationSettings{
		EmailEnabled:    true,
		BedtimeReminder: true,
		ReminderTime:    "22:00",
		WeeklyReport:    true,
	}

	data := &TemplateData{
		Title:      "設定",
		ActiveMenu: "settings",
		User:       user,
		Data: map[string]interface{}{
			"Notifications": notifications,
		},
	}

	err := h.templates.Render(w, "settings.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	UpdateProfile プロフィールの更新
*/
func (h *SettingsHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// TODO: プロフィール更新の実装
	// name := r.FormValue("name")
	// email := r.FormValue("email")
	// timezone := r.FormValue("timezone")

	data := &TemplateData{
		Title:      "設定",
		ActiveMenu: "settings",
		Flash: &Flash{
			Type:    "success",
			Message: "プロフィールを更新しました",
		},
	}

	err := h.templates.Render(w, "settings.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	UpdatePassword パスワードの更新
*/
func (h *SettingsHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	// パスワードが変更されている？
	if currentPassword == newPassword {
		data := &TemplateData{
			Title:      "設定",
			ActiveMenu: "settings",
			Flash: &Flash{
				Type:    "danger",
				Message: "新しいパスワードが古いパスワードと同じです",
			},
		}
		h.templates.Render(w, "reset-password.html", data)
		return
	}

	// パスワードのバリデーション
	if newPassword != confirmPassword {
		data := &TemplateData{
			Title:      "設定",
			ActiveMenu: "settings",
			Flash: &Flash{
				Type:    "danger",
				Message: "新しいパスワードが一致しません",
			},
		}
		h.templates.Render(w, "reset-password.html", data)
		return
	}

	// パスワードの条件が満たされているかどうか
	passwordPattern := `^[A-Za-z0-9!@#$%^&*()_+\-=\[\]{};':\"\\|,.<>\/?]{8,}$`
	passwordPatternRegex := regexp.MustCompile(passwordPattern)
	if !passwordPatternRegex.MatchString(newPassword) {
		data := &TemplateData{
			Title: "新しいパスワードの設定",
			Flash: &Flash{
				Type:    "danger",
				Message: "パスワードは半角英数字と記号のみ使用できます。",
			},
		}
		h.templates.Render(w, "reset-password.html", data)
		return
	}

	// TODO: パスワード更新の実装

	data := &TemplateData{
		Title:      "設定",
		ActiveMenu: "settings",
		Flash: &Flash{
			Type:    "success",
			Message: "パスワードを更新しました",
		},
	}

	err := h.templates.Render(w, "settings.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	UpdateNotifications 通知設定の更新
*/
func (h *SettingsHandler) UpdateNotifications(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// TODO: 通知設定の更新の実装
	// emailEnabled := r.FormValue("email_notification") == "on"
	// bedtimeReminder := r.FormValue("bedtime_reminder") == "on"
	// reminderTime := r.FormValue("reminder_time")
	// weeklyReport := r.FormValue("weekly_report") == "on"

	data := &TemplateData{
		Title:      "設定",
		ActiveMenu: "settings",
		Flash: &Flash{
			Type:    "success",
			Message: "通知設定を更新しました",
		},
	}

	err := h.templates.Render(w, "settings.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	ShowDeleteAccountPage アカウント削除の確認
*/
func (h *SettingsHandler) ShowDeleteAccountPage(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:      "アカウント削除の確認",
		ActiveMenu: "settings",
	}

	err := h.templates.Render(w, "delete-account.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	DeleteAccount アカウントの削除
*/
func (h *SettingsHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// password := r.FormValue("password")
	// confirmDelete := r.FormValue("confirm_delete") == "on"

	// TODO: アカウント削除の実装

	// セッションをクリアしてログインページにリダイレクト
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

/*
	ExportCSV 睡眠記録のCSVエクスポート
*/
func (h *SettingsHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	// ファイル名の設定
	filename := fmt.Sprintf("sleep-records-%s.csv", time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// CSVライターの作成
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// ヘッダーの書き込み
	headers := []string{"日付", "就寝時刻", "起床時刻", "睡眠時間", "睡眠スコア", "睡眠の質", "メモ"}
	if err := writer.Write(headers); err != nil {
		http.Error(w, "CSVの書き込みに失敗しました", http.StatusInternalServerError)
		return
	}

	// TODO: データの書き込み
	// データベースからユーザーの睡眠記録を取得して書き込む
}

/*
	ExportJSON 睡眠記録のJSONエクスポート
*/
func (h *SettingsHandler) ExportJSON(w http.ResponseWriter, r *http.Request) {
	// ファイル名の設定
	filename := fmt.Sprintf("sleep-records-%s.json", time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// TODO: データのエクスポート
	// JSONエクスポート用の構造体
	type SleepRecordExport struct {
		Date     string  `json:"date"`
		BedTime  string  `json:"bed_time"`
		WakeTime string  `json:"wake_time"`
		Duration float64 `json:"duration"`
		Score    int     `json:"score"`
		Quality  string  `json:"quality"`
		Notes    string  `json:"notes"`
	}

	data := struct {
		ExportedAt time.Time          `json:"exported_at"`
		Records    []SleepRecordExport `json:"records"`
	}{
		ExportedAt: time.Now(),
		Records:    make([]SleepRecordExport, 0), // TODO: データベースからデータを取得
	}

	// JSONエンコード
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		http.Error(w, "JSONの書き込みに失敗しました", http.StatusInternalServerError)
		return
	}
}

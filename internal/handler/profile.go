// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/profile.go
// profileは、プロフィール関連のハンドラーを提供します。

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
	"github.com/223n-tech/SuiminNisshi-Go/internal/util"
)

// プロフィール関連のハンドラー
type ProfileHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// プロフィール更新リクエストの構造体
type ProfileUpdateRequest struct {
	DisplayName        string `json:"display_name"`
	Email             string `json:"email"`
	CurrentPassword   string `json:"current_password"`
	NewPassword       string `json:"new_password"`
	PasswordConfirm   string `json:"password_confirm"`
	Timezone         string `json:"timezone"`
}

// ProfileHandlerを作成
func NewProfileHandler(templates *TemplateManager, svc *service.Service) *ProfileHandler {
	return &ProfileHandler{
		templates: templates,
		service:   svc,
	}
}

// ルーティングを登録
func (h *ProfileHandler) RegisterRoutes(r *RouterWrapper) {
	r.Get("/profile", h.Profile)
	r.Post("/profile/update", h.UpdateProfile)
	r.Post("/profile/password", h.UpdatePassword)
	r.Post("/profile/preferences", h.UpdatePreferences)
}

// プロフィール画面を表示
func (h *ProfileHandler) Profile(w http.ResponseWriter, r *http.Request) {
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
	preferences, err := h.service.User().GetSleepPreference(r.Context(), userID)
	if err != nil {
		http.Error(w, "睡眠設定の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// 利用可能なタイムゾーンのリスト
	timezones := []string{
		"Asia/Tokyo",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
		// 他のタイムゾーン...
	}

	data := &TemplateData{
		Title:      "プロフィール",
		ActiveMenu: "profile",
		User:       user,
		Data: map[string]interface{}{
			"Preferences": preferences,
			"Timezones":   timezones,
		},
	}

	err = h.templates.Render(w, "profile.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// プロフィールの更新
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// リクエストデータの取得
	email := r.FormValue("email")
	displayName := r.FormValue("display_name")
	timezone := r.FormValue("timezone")

	// バリデーション
	if err := h.validateProfileData(displayName, email, timezone); err != nil {
		data := &TemplateData{
			Title:      "プロフィール",
			ActiveMenu: "profile",
			Flash: &Flash{
				Type:    "danger",
				Message: err.Error(),
			},
		}
		h.templates.Render(w, "profile.html", data)
		return
	}

	// ユーザー情報の更新
	user := &models.User{
		ID:          userID,
		Email:       email,
		DisplayName: displayName,
	}

	err := h.service.User().UpdateProfile(r.Context(), user)
	if err != nil {
		data := &TemplateData{
			Title:      "プロフィール",
			ActiveMenu: "profile",
			Flash: &Flash{
				Type:    "danger",
				Message: "プロフィールの更新に失敗しました",
			},
		}
		h.templates.Render(w, "profile.html", data)
		return
	}

	// 成功メッセージを表示して再表示
	data := &TemplateData{
		Title:      "プロフィール",
		ActiveMenu: "profile",
		User:       user,
		Flash: &Flash{
			Type:    "success",
			Message: "プロフィールを更新しました",
		},
	}
	h.templates.Render(w, "profile.html", data)
}

// パスワードの更新
func (h *ProfileHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("password_confirm")

	if (newPassword != confirmPassword) {
		data := &TemplateData{
			Title:      "プロフィール",
			ActiveMenu: "profile",
			Flash: &Flash{
				Type:    "danger",
				Message: "新しいパスワードが一致しません",
			},
		}
		h.templates.Render(w, "profile.html", data)
		return
	}

	// パスワードの更新処理
	err := h.service.User().UpdatePassword(r.Context(), userID, currentPassword, newPassword)
	if err != nil {
		data := &TemplateData{
			Title:      "プロフィール",
			ActiveMenu: "profile",
			Flash: &Flash{
				Type:    "danger",
				Message: "パスワードの更新に失敗しました: " + err.Error(),
			},
		}
		h.templates.Render(w, "profile.html", data)
		return
	}

	// 成功メッセージを表示して再表示
	data := &TemplateData{
		Title:      "プロフィール",
		ActiveMenu: "profile",
		Flash: &Flash{
			Type:    "success",
			Message: "パスワードを更新しました",
		},
	}
	h.templates.Render(w, "profile.html", data)
}

// 睡眠設定の更新
func (h *ProfileHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// 設定データの取得と更新
	pref := &models.UserSleepPreference{
		UserID:              userID,
		PreferredBedtime:    util.ParseTime(r.FormValue("preferred_bedtime")),
		PreferredWakeupTime: util.ParseTime(r.FormValue("preferred_wakeup_time")),
		SleepGoalHours:      parseInt(r.FormValue("sleep_goal_hours"), 8),
		IsReminderEnabled:   r.FormValue("is_reminder_enabled") == "on",
	}

	err := h.service.User().UpdateSleepPreference(r.Context(), pref)
	if err != nil {
		data := &TemplateData{
			Title:      "プロフィール",
			ActiveMenu: "profile",
			Flash: &Flash{
				Type:    "danger",
				Message: "睡眠設定の更新に失敗しました",
			},
		}
		h.templates.Render(w, "profile.html", data)
		return
	}

	// 成功メッセージを表示して再表示
	data := &TemplateData{
		Title:      "プロフィール",
		ActiveMenu: "profile",
		Flash: &Flash{
			Type:    "success",
			Message: "睡眠設定を更新しました",
		},
	}
	h.templates.Render(w, "profile.html", data)
}

// プロフィールデータのバリデーション
func (h *ProfileHandler) validateProfileData(displayName, email, timezone string) error {
	if displayName == "" {
		return service.ErrEmptyDisplayName
	}

	if email == "" {
		return service.ErrEmptyEmail
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return service.ErrInvalidEmail
	}

	if timezone == "" {
		return service.ErrEmptyTimezone
	}

	// TODO: タイムゾーンの有効性チェック
	return nil
}

// 文字列を整数に変換
func parseInt(str string, defaultValue int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return value
}

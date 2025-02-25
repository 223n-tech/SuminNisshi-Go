// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/password_reset.go
// password_resetは、パスワードリセット関連のハンドラーを提供します。

import (
	"net/http"
	"regexp"

	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
	"github.com/go-chi/chi/v5"
)

// パスワードリセット関連のハンドラー
type PasswordResetHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// PasswordResetHandlerを作成
func NewPasswordResetHandler(templates *TemplateManager, svc *service.Service) *PasswordResetHandler {
	return &PasswordResetHandler{
		templates: templates,
		service:   svc,
	}
}

// ルーティングを登録
func (h *PasswordResetHandler) RegisterRoutes(r chi.Router) {
	r.Get("/forgot-password", h.ShowForgotPasswordForm)
	r.Post("/forgot-password", h.HandleForgotPassword)
	r.Get("/reset-password/{token}", h.ShowResetPasswordForm)
	r.Post("/reset-password/{token}", h.HandleResetPassword)
}

// パスワードリセットフォームの表示
func (h *PasswordResetHandler) ShowForgotPasswordForm(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "パスワードの再設定",
	}

	// フラッシュメッセージの取得
	if flashMsg := r.URL.Query().Get("message"); flashMsg != "" {
		data.Flash = &Flash{
			Type:    r.URL.Query().Get("type"),
			Message: flashMsg,
		}
	}

	err := h.templates.Render(w, "forgot-password.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// パスワードリセットリクエストの処理
func (h *PasswordResetHandler) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		data := &TemplateData{
			Title: "パスワードの再設定",
			Flash: &Flash{
				Type:    "danger",
				Message: "メールアドレスを入力してください",
			},
		}
		h.templates.Render(w, "forgot-password.html", data)
		return
	}

	// メールアドレスのバリデーション
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		data := &TemplateData{
			Title: "パスワードの再設定",
			Flash: &Flash{
				Type:    "danger",
				Message: "有効なメールアドレスを入力してください",
			},
		}
		h.templates.Render(w, "forgot-password.html", data)
		return
	}

	// パスワードリセットトークンの生成と保存
	_, err := h.service.User().InitiatePasswordReset(r.Context(), email)
	if err != nil {
		h.service.Logger().Printf("[ERROR] パスワードリセットの初期化に失敗: error=%v, email=%s", err, email)
		data := &TemplateData{
			Title: "パスワードの再設定",
			Flash: &Flash{
				Type:    "danger",
				Message: "パスワードリセットの処理中にエラーが発生しました",
			},
		}
		h.templates.Render(w, "forgot-password.html", data)
		return
	}

	// 成功メッセージを表示
	data := &TemplateData{
		Title: "パスワードの再設定",
		Flash: &Flash{
			Type:    "success",
			Message: "パスワード再設定用のメールを送信しました。メールの指示に従って操作を完了してください。",
		},
	}

	h.templates.Render(w, "forgot-password.html", data)
}

// パスワード再設定フォームの表示
func (h *PasswordResetHandler) ShowResetPasswordForm(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}

	// トークンの検証
	valid, err := h.service.User().ValidateResetToken(r.Context(), token)
	if err != nil || !valid {
		data := &TemplateData{
			Title: "無効なリンク",
			Flash: &Flash{
				Type:    "danger",
				Message: "パスワードリセットリンクが無効か期限切れです。再度パスワードリセットを実行してください。",
			},
		}
		h.templates.Render(w, "forgot-password.html", data)
		return
	}

	data := &TemplateData{
		Title: "新しいパスワードの設定",
		Data: map[string]interface{}{
			"Token": token,
		},
	}

	err = h.templates.Render(w, "reset-password.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// パスワード再設定の処理
func (h *PasswordResetHandler) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirmation")

	// パスワードの検証
	if err := h.validatePassword(password, passwordConfirm); err != nil {
		data := &TemplateData{
			Title: "新しいパスワードの設定",
			Flash: &Flash{
				Type:    "danger",
				Message: err.Error(),
			},
			Data: map[string]interface{}{
				"Token": token,
			},
		}
		h.templates.Render(w, "reset-password.html", data)
		return
	}

	// パスワードの更新
	err := h.service.User().CompletePasswordReset(r.Context(), token, password)
	if err != nil {
		data := &TemplateData{
			Title: "新しいパスワードの設定",
			Flash: &Flash{
				Type:    "danger",
				Message: "パスワードの更新に失敗しました",
			},
			Data: map[string]interface{}{
				"Token": token,
			},
		}
		h.templates.Render(w, "reset-password.html", data)
		return
	}

	// ログインページにリダイレクト（成功メッセージ付き）
	http.Redirect(w, r, "/login?message=パスワードが正常に更新されました&type=success", http.StatusSeeOther)
}

// パスワードのバリデーション
func (h *PasswordResetHandler) validatePassword(password, confirm string) error {
	if password == "" {
		return service.ErrEmptyPassword
	}

	if password != confirm {
		return service.ErrPasswordMismatch
	}

	if len(password) < 8 {
		return service.ErrPasswordTooShort
	}

	passwordPattern := `^[A-Za-z0-9!@#$%^&*()_+\-=\[\]{};':\"\\|,.<>\/?]{8,}$`
	passwordPatternRegex := regexp.MustCompile(passwordPattern)
	if !passwordPatternRegex.MatchString(password) {
		return service.ErrInvalidPasswordFormat
	}

	return nil
}

package handler

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
)

// PasswordResetHandler パスワードリセット関連のハンドラー
type PasswordResetHandler struct {
	templates *TemplateManager
}

// NewPasswordResetHandler パスワードリセットハンドラーを作成
func NewPasswordResetHandler(templates *TemplateManager) *PasswordResetHandler {
	return &PasswordResetHandler{
		templates: templates,
	}
}

// RegisterRoutes ルートの登録
func (h *PasswordResetHandler) RegisterRoutes(r chi.Router) {
	r.Get("/forgot-password", h.ShowForgotPasswordForm)
	r.Post("/forgot-password", h.HandleForgotPassword)
	r.Get("/reset-password/{token}", h.ShowResetPasswordForm)
	r.Post("/reset-password/{token}", h.HandleResetPassword)
}

// ShowForgotPasswordForm パスワード忘れフォームの表示
func (h *PasswordResetHandler) ShowForgotPasswordForm(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "パスワードの再設定",
	}

	err := h.templates.Render(w, "forgot-password.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleForgotPassword パスワードリセットメールの送信処理
func (h *PasswordResetHandler) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: メール送信処理の実装
	_ = r.FormValue("email")

	// TODO: メール送信処理の実装
	// 1. メールアドレスの検証
	// 2. トークンの生成
	// 3. トークンの保存
	// 4. メール送信

	data := &TemplateData{
		Title: "パスワードの再設定",
		Flash: &Flash{
			Type:    "success",
			Message: "パスワード再設定用のメールを送信しました。メールの指示に従って操作を完了してください。",
		},
	}

	err := h.templates.Render(w, "forgot-password.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ShowResetPasswordForm パスワード再設定フォームの表示
func (h *PasswordResetHandler) ShowResetPasswordForm(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	// TODO: トークンの検証
	// 1. トークンの有効性チェック
	// 2. トークンの有効期限チェック

	data := &TemplateData{
		Title: "新しいパスワードの設定",
		Data: map[string]interface{}{
			"Token": token,
		},
	}

	err := h.templates.Render(w, "reset-password.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleResetPassword パスワード再設定の処理
func (h *PasswordResetHandler) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirmation")

	// パスワードが一致しているかどうか
	if password != passwordConfirm {
		data := &TemplateData{
			Title: "新しいパスワードの設定",
			Flash: &Flash{
				Type:    "danger",
				Message: "パスワードが一致しません。",
			},
			Data: map[string]interface{}{
				"Token": token,
			},
		}
		h.templates.Render(w, "reset-password.html", data)
		return
	}

	// パスワードの条件が満たされているかどうか
	passwordPattern := `^[A-Za-z0-9!@#$%^&*()_+\-=\[\]{};':\"\\|,.<>\/?]{8,}$`
	passwordPatternRegex := regexp.MustCompile(passwordPattern)
	if !passwordPatternRegex.MatchString(password) {
		data := &TemplateData{
			Title: "新しいパスワードの設定",
			Flash: &Flash{
				Type:    "danger",
				Message: "パスワードは半角英数字と記号のみ使用できます。",
			},
			Data: map[string]interface{}{
				"Token": token,
			},
		}
		h.templates.Render(w, "reset-password.html", data)
		return
	}

	// TODO: パスワード更新処理の実装
	// 1. トークンの検証
	// 2. パスワードのハッシュ化
	// 3. パスワードの更新
	// 4. トークンの無効化

	// 成功時はログインページにリダイレクト
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

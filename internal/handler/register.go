// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/register.go
// registerは、新規登録関連のハンドラーを提供します。

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
	"github.com/go-chi/chi/v5"
)

// 新規登録画面のハンドラー
type RegisterHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// 新規登録画面のデータ
type RegisterData struct {
	Name                 string
	Email                string
	Password             string
	PasswordConfirmation string
	Terms                bool
	Error                string
}

// RegisterHandlerを作成
func NewRegisterHandler(templates *TemplateManager, svc *service.Service) *RegisterHandler {
	return &RegisterHandler{
		templates: templates,
		service:   svc,
	}
}

// ルーティングを登録
func (h *RegisterHandler) RegisterRoutes(r chi.Router) {
	r.Get("/register", h.RegisterPage)
	r.Post("/register", h.Register)
}

// 新規登録画面を表示
func (h *RegisterHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	// すでにログインしている場合はダッシュボードにリダイレクト
	if user := r.Context().Value(UserKey); user != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	data := &TemplateData{
		Title: "アカウント登録",
		Data: map[string]interface{}{
			"Form": &RegisterData{},
		},
	}

	err := h.templates.Render(w, "register.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 新規登録処理
func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// フォームデータの取得
	registerData := &RegisterData{
		Name:                 r.FormValue("name"),
		Email:                r.FormValue("email"),
		Password:             r.FormValue("password"),
		PasswordConfirmation: r.FormValue("password_confirmation"),
		Terms:                r.FormValue("terms") == "on",
	}

	// バリデーション
	if err := h.validateRegisterData(registerData); err != nil {
		data := &TemplateData{
			Title: "アカウント登録",
			Data: map[string]interface{}{
				"Form": registerData,
			},
			Flash: &Flash{
				Type:    "danger",
				Message: err.Error(),
			},
		}
		h.templates.Render(w, "register.html", data)
		return
	}

	// ユーザー登録
	_, err := h.service.User().Register(r.Context(), registerData.Email, registerData.Name, registerData.Password)
	if err != nil {
		// 既存のメールアドレスの場合
		if err == service.ErrEmailAlreadyExists {
			data := &TemplateData{
				Title: "アカウント登録",
				Data: map[string]interface{}{
					"Form": registerData,
				},
				Flash: &Flash{
					Type:    "danger",
					Message: "このメールアドレスは既に登録されています",
				},
			}
			h.templates.Render(w, "register.html", data)
			return
		}

		// その他のエラー
		h.service.Logger().Printf("[ERROR] ユーザー登録エラー: error=%v", err)
		data := &TemplateData{
			Title: "アカウント登録",
			Data: map[string]interface{}{
				"Form": registerData,
			},
			Flash: &Flash{
				Type:    "danger",
				Message: "ユーザー登録に失敗しました",
			},
		}
		h.templates.Render(w, "register.html", data)
		return
	}

	// 確認メール送信
	err = h.service.Email().SendWelcomeEmail(r.Context(), registerData.Email, registerData.Name)
	if err != nil {
		h.service.Logger().Printf("[ERROR] 確認メール送信エラー: error=%v", err)
	}

	// 登録成功時の処理
	http.Redirect(w, r, "/register/complete", http.StatusSeeOther)
}

// 新規登録データのバリデーション
func (h *RegisterHandler) validateRegisterData(data *RegisterData) error {
	if data.Name == "" {
		return fmt.Errorf("名前を入力してください")
	}

	if data.Email == "" {
		return fmt.Errorf("メールアドレスを入力してください")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(data.Email) {
		return fmt.Errorf("有効なメールアドレスを入力してください")
	}

	if data.Password == "" {
		return fmt.Errorf("パスワードを入力してください")
	}

	if len(data.Password) < 8 {
		return fmt.Errorf("パスワードは8文字以上である必要があります")
	}

	passwordPattern := `^[A-Za-z0-9!@#$%^&*()_+\-=\[\]{};':\"\\|,.<>\/?]{8,}$`
	passwordPatternRegex := regexp.MustCompile(passwordPattern)
	if !passwordPatternRegex.MatchString(data.Password) {
		return fmt.Errorf("パスワードは半角英数字と記号のみ使用できます")
	}

	if data.Password != data.PasswordConfirmation {
		return fmt.Errorf("パスワードが一致しません")
	}

	if !data.Terms {
		return fmt.Errorf("利用規約に同意する必要があります")
	}

	return nil
}

// 登録完了画面を表示
func (h *RegisterHandler) RegistrationComplete(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "登録完了",
		Flash: &Flash{
			Type:    "success",
			Message: "アカウント登録が完了しました。メールをご確認ください。",
		},
	}

	err := h.templates.Render(w, "register-complete.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

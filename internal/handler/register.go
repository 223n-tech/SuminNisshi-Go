package handler

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
)

// RegisterHandler 新規登録画面のハンドラー
type RegisterHandler struct {
	templates *TemplateManager
}

// RegisterData 登録フォームのデータ
type RegisterData struct {
	Name                 string
	Email                string
	Password             string
	PasswordConfirmation string
	Terms                bool
	Error               string
}

// NewRegisterHandler 新規登録ハンドラーを作成
func NewRegisterHandler(templates *TemplateManager) *RegisterHandler {
	return &RegisterHandler{
		templates: templates,
	}
}

// RegisterRoutes ルートの登録
func (h *RegisterHandler) RegisterRoutes(r chi.Router) {
	r.Get("/register", h.RegisterPage)
	r.Post("/register", h.Register)
}

// RegisterPage 新規登録ページの表示
func (h *RegisterHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "アカウント登録",
	}

	err := h.templates.Render(w, "register.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Register 新規登録の処理
func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	data := &RegisterData{
		Name:                 r.FormValue("name"),
		Email:                r.FormValue("email"),
		Password:             r.FormValue("password"),
		PasswordConfirmation: r.FormValue("password_confirmation"),
		Terms:                r.FormValue("terms") == "on",
	}

	// バリデーション
	if err := h.validateRegisterData(data); err != nil {
		data := &TemplateData{
			Title: "アカウント登録",
			Data:  data,
			Flash: &Flash{
				Type:    "danger",
				Message: err.Error(),
			},
		}
		h.templates.Render(w, "register.html", data)
		return
	}

	// TODO: ユーザー登録処理の実装

	// 登録成功後、ログインページにリダイレクト
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// validateRegisterData 登録データのバリデーション
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

	if data.Password != data.PasswordConfirmation {
		return fmt.Errorf("パスワードが一致しません")
	}

	if !data.Terms {
		return fmt.Errorf("利用規約に同意する必要があります")
	}

	return nil
}

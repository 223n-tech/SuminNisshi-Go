// 認証関連のハンドラを定義する
package handler

import (
	"net/http"
)

/*
	AuthHandler 認証関連のハンドラ
*/
type AuthHandler struct {
	templates *TemplateManager
}

// NewAuthHandler は AuthHandler を作成します。
func NewAuthHandler(templates *TemplateManager) *AuthHandler {
	return &AuthHandler{
		templates: templates,
	}
}

/*
	RegisterRoutes ルーティングを登録
*/
func (h *AuthHandler) RegisterRoutes(r *RouterWrapper) {
	r.Get("/login", h.LoginPage)
	r.Post("/login", h.Login)
	r.Get("/logout", h.Logout)
}

/*
	LoginPage ログイン画面を表示
*/
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := h.templates.Render(w, "login.html", &TemplateData{
		Title: "ログイン",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	Login ログイン処理
*/
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: 実際のログイン処理を実装
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

/*
	Logout ログアウト処理
*/
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: ログアウト処理を実装
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

/*
	RequireAuth 認証が必要なリクエストに対してミドルウェアを適用
*/
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: セッションまたはJWTトークンの検証
		// とりあえず開発用に常にtrueを返す
		next.ServeHTTP(w, r)
	})
}

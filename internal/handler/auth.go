package handler

import (
	"net/http"
)

type AuthHandler struct {
	templates *TemplateManager
}

func NewAuthHandler(templates *TemplateManager) *AuthHandler {
	return &AuthHandler{
		templates: templates,
	}
}

func (h *AuthHandler) RegisterRoutes(r *RouterWrapper) {
	r.Get("/login", h.LoginPage)
	r.Post("/login", h.Login)
	r.Get("/logout", h.Logout)
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := h.templates.Render(w, "login.html", &TemplateData{
		Title: "ログイン",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: 実際のログイン処理を実装
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: ログアウト処理を実装
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

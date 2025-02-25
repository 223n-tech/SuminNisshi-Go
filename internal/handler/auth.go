// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/auth.go
// authは、認証関連のハンドラーを提供します。

import (
	"context"
	"net/http"

	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
)

// コンテキストのキー
type userContextKey string

// ユーザー情報のコンテキストキー
const UserKey userContextKey = "user"

// 認証関連のハンドラー
type AuthHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// AuthHandlerを作成
func NewAuthHandler(templates *TemplateManager, svc *service.Service) *AuthHandler {
	return &AuthHandler{
		templates: templates,
		service:   svc,
	}
}

// ルーティングを登録
func (h *AuthHandler) RegisterRoutes(r *RouterWrapper) {
	r.Get("/login", h.LoginPage)
	r.Post("/login", h.Login)
	r.Get("/logout", h.Logout)
}

// ログイン画面を表示
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// すでにログインしている場合はダッシュボードにリダイレクト
	if user := r.Context().Value(UserKey); user != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	data := &TemplateData{
		Title: "ログイン",
	}

	// フラッシュメッセージの取得
	if flashMsg := r.URL.Query().Get("message"); flashMsg != "" {
		data.Flash = &Flash{
			Type:    r.URL.Query().Get("type"),
			Message: flashMsg,
		}
	}

	err := h.templates.Render(w, "login.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ログイン処理
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// 入力値の検証
	if email == "" || password == "" {
		data := &TemplateData{
			Title: "ログイン",
			Flash: &Flash{
				Type:    "danger",
				Message: "メールアドレスとパスワードを入力してください",
			},
		}
		h.templates.Render(w, "login.html", data)
		return
	}

	// ユーザー認証
	_, err := h.service.User().Authenticate(r.Context(), email, password)
	if err != nil {
		data := &TemplateData{
			Title: "ログイン",
			Flash: &Flash{
				Type:    "danger",
				Message: "メールアドレスまたはパスワードが正しくありません",
			},
		}
		h.templates.Render(w, "login.html", data)
		return
	}

	// セッションの作成
	// TODO: セッション管理の実装
	// session := CreateSession(r.Context(), user)

	// TODO: セッションIDをクッキーに設定
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		// Value:    session.ID,
		Value:    "dummy-session-id",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// ダッシュボードにリダイレクト
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// ログアウト処理
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// セッションの破棄
	// TODO: セッション管理の実装
	// ClearSession(r.Context())

	// セッションクッキーの削除
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// ログインページにリダイレクト
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// 認証が必要なリクエストに対してミドルウェアを適用
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: セッションの検証
		// sessionID, err := r.Cookie("session_id")
		// if err != nil {
		// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
		// 	return
		// }

		// TODO: セッションからユーザー情報を取得
		// user, err := GetUserFromSession(r.Context(), sessionID.Value)
		// if err != nil {
		// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
		// 	return
		// }

		// ユーザー情報をコンテキストに設定
		// ctx := context.WithValue(r.Context(), UserKey, user)
		// next.ServeHTTP(w, r.WithContext(ctx))

		// 開発用に一時的にスキップ
		next.ServeHTTP(w, r)
	})
}

// コンテキストからユーザー情報を取得
func GetUserFromContext(ctx context.Context) interface{} {
	return ctx.Value(UserKey)
}

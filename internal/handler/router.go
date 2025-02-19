package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RouterWrapper はchi.Routerのラッパー構造体です
type RouterWrapper struct {
	chi.Router
}

// NewRouter は新しいRouterWrapperインスタンスを作成します
func NewRouter(r chi.Router) *RouterWrapper {
	return &RouterWrapper{Router: r}
}

// Group はルートグループを作成します
func (r *RouterWrapper) Group(pattern string, fn func(r chi.Router)) chi.Router {
	return r.Router.Group(fn)
}

// Static は静的ファイルを提供するルートを設定します
func (r *RouterWrapper) Static(prefix string, root http.FileSystem) {
	r.Router.Handle(prefix+"/*", http.StripPrefix(prefix, http.FileServer(root)))
}

// SubRouter は新しいサブルーターを作成します
func (r *RouterWrapper) SubRouter(prefix string) *RouterWrapper {
	return NewRouter(chi.NewRouter())
}

// WithMiddleware はミドルウェア付きの新しいルーターを作成します
func (r *RouterWrapper) WithMiddleware(middlewares ...func(http.Handler) http.Handler) *RouterWrapper {
	router := chi.NewRouter()
	router.Use(middlewares...)
	return NewRouter(router)
}

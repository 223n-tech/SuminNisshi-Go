// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/router.go
// routerは、ルーターのラッパー構造体を提供します

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

/*
	chi.Routerのラッパーです
*/
type RouterWrapper struct {
	chi.Router
}

/*
	新しいRouterWrapperを作成します
*/
func NewRouter(r chi.Router) *RouterWrapper {
	return &RouterWrapper{Router: r}
}

/*
	グループを作成します
*/
func (r *RouterWrapper) Group(_ string, fn func(r chi.Router)) chi.Router {
	return r.Router.Group(fn)
}

/*
	静的ファイルを提供します
*/
func (r *RouterWrapper) Static(prefix string, root http.FileSystem) {
	r.Router.Handle(prefix+"/*", http.StripPrefix(prefix, http.FileServer(root)))
}

/*
	サブルーターを作成します
*/
func (r *RouterWrapper) SubRouter(_ string) *RouterWrapper {
	return NewRouter(chi.NewRouter())
}

/*
	指定されたミドルウェアを適用した新しいRouterWrapperを返します
*/
func (r *RouterWrapper) WithMiddleware(middlewares ...func(http.Handler) http.Handler) *RouterWrapper {
	router := chi.NewRouter()
	router.Use(middlewares...)
	return NewRouter(router)
}

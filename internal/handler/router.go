// ルーターのラッパー構造体を提供します
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

/*
	RouterWrapper はchi.Routerのラッパーです
*/
type RouterWrapper struct {
	chi.Router
}

/*
	NewRouter は新しいRouterWrapperを作成します
*/
func NewRouter(r chi.Router) *RouterWrapper {
	return &RouterWrapper{Router: r}
}

/*
	Get はGETリクエストを処理するハンドラを登録します
*/
func (r *RouterWrapper) Group(pattern string, fn func(r chi.Router)) chi.Router {
	return r.Router.Group(fn)
}

/*
	Get はGETリクエストを処理するハンドラを登録します
*/
func (r *RouterWrapper) Static(prefix string, root http.FileSystem) {
	r.Router.Handle(prefix+"/*", http.StripPrefix(prefix, http.FileServer(root)))
}

/*
	Get はGETリクエストを処理するハンドラを登録します
*/
func (r *RouterWrapper) SubRouter(prefix string) *RouterWrapper {
	return NewRouter(chi.NewRouter())
}

/*
	WithMiddleware は指定されたミドルウェアを適用した新しいRouterWrapperを返します
*/
func (r *RouterWrapper) WithMiddleware(middlewares ...func(http.Handler) http.Handler) *RouterWrapper {
	router := chi.NewRouter()
	router.Use(middlewares...)
	return NewRouter(router)
}

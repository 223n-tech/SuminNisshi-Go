package handler

import (
	"log"
	"net/http"

	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
)

// DashboardHandler ダッシュボード関連のハンドラー
type DashboardHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// NewDashboardHandler ダッシュボードハンドラーを作成
func NewDashboardHandler(templates *TemplateManager, svc *service.Service) *DashboardHandler {
	return &DashboardHandler{
		templates: templates,
		service:   svc,
	}
}

// RegisterRoutes ルートの登録
func (h *DashboardHandler) RegisterRoutes(r *RouterWrapper) {
	r.Get("/dashboard", h.Dashboard)
}

// Dashboard ダッシュボードページの表示
func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:      "ダッシュボード",
		ActiveMenu: "dashboard",
	}

	// テンプレートのデバッグ情報を出力
	log.Printf("Rendering dashboard template with data: %+v", data)
	err := h.templates.Render(w, "dashboard.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/dashboard.go
// dashboardは、ダッシュボード関連のハンドラーを提供します。

import (
	"net/http"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
)

// ダッシュボード関連のハンドラー
type DashboardHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// DashboardHandlerを作成
func NewDashboardHandler(templates *TemplateManager, svc *service.Service) *DashboardHandler {
	return &DashboardHandler{
		templates: templates,
		service:   svc,
	}
}

// ルーティングを登録
func (h *DashboardHandler) RegisterRoutes(r *RouterWrapper) {
	r.Get("/dashboard", h.Dashboard)
	r.Get("/api/dashboard/summary", h.GetDashboardSummary)
}

// ダッシュボードを表示
func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// ユーザー情報の取得
	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userModel, ok := user.(*models.User)
	if !ok {
		// TODO: 実装
	}

	// TODO: 実際のユーザーIDを使用
	// userID := user.ID
	var userID int64 = 1 // 開発用

	// 睡眠設定の取得
	pref, err := h.service.User().GetSleepPreference(r.Context(), userID)
	if err != nil {
		h.templates.Render(w, "500.html", &TemplateData{
			Title: "Internal Server Error",
		})
		return
	}

	// 直近の睡眠記録を取得
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)
	diary, err := h.service.Diary().GetDiaryByDateRange(r.Context(), userID, startDate, endDate)
	if err != nil {
		h.templates.Render(w, "500.html", &TemplateData{
			Title: "Internal Server Error",
		})
		return
	}

	// 統計データの取得
	var stats *service.DashboardStats
	if len(diary) > 0 {
		stats, err = h.service.Record().GetDashboardStats(r.Context(), diary[0].ID)
		if err != nil {
			h.templates.Render(w, "500.html", &TemplateData{
				Title: "Internal Server Error",
			})
			return
		}
	}

	data := &TemplateData{
		Title:      "ダッシュボード",
		ActiveMenu: "dashboard",
		User:       userModel,
		Data: map[string]interface{}{
			"Preferences": pref,
			"Statistics": stats,
			"StartDate":  startDate.Format("2006-01-02"),
			"EndDate":    endDate.Format("2006-01-02"),
		},
	}

	err = h.templates.Render(w, "dashboard.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ダッシュボードのサマリーデータを取得
func (h *DashboardHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// 期間の取得
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)

	// 睡眠記録の取得
	diary, err := h.service.Diary().GetDiaryByDateRange(r.Context(), userID, startDate, endDate)
	if err != nil {
		http.Error(w, "睡眠記録の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	if len(diary) == 0 {
		http.Error(w, "睡眠記録が見つかりません", http.StatusNotFound)
		return
	}

	// 統計データの取得
	stats, err := h.service.Record().GetDashboardStats(r.Context(), diary[0].ID)
	if err != nil {
		http.Error(w, "統計データの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// JSONでレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.Write(stats.ToJSON())
}


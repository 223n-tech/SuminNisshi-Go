// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/statistics.go
// statisticsは、統計情報画面のハンドラーを実装しています。
// 統計情報画面は、睡眠記録の統計情報を表示する画面です。
// 統計情報は、睡眠時間やスコアの平均値、スコアの分布、月間のサマリーなどを表示します。

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
	"github.com/go-chi/chi/v5"
)

// 統計情報関連のハンドラー
type StatisticsHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// StatisticsHandlerを作成
func NewStatisticsHandler(templates *TemplateManager, svc *service.Service) *StatisticsHandler {
	return &StatisticsHandler{
		templates: templates,
		service:   svc,
	}
}

// ルーティングを登録
func (h *StatisticsHandler) RegisterRoutes(r chi.Router) {
	r.Get("/statistics", h.Statistics)
	r.Get("/api/statistics/data", h.GetStatisticsData)
	r.Get("/api/statistics/weekly", h.GetWeeklyStats)
	r.Get("/api/statistics/monthly", h.GetMonthlyStats)
}

// 統計情報画面を表示
func (h *StatisticsHandler) Statistics(w http.ResponseWriter, r *http.Request) {
	// ユーザー情報の取得
	userCtx := r.Context().Value(UserKey)
	if userCtx == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// デフォルトの期間を設定（直近30日）
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	// 睡眠設定の取得
	pref, err := h.service.User().GetSleepPreference(r.Context(), userID)
	if err != nil {
		http.Error(w, "睡眠設定の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	data := &TemplateData{
		Title:      "統計情報",
		ActiveMenu: "statistics",
		Data: map[string]interface{}{
			"StartDate":    startDate.Format("2006-01-02"),
			"EndDate":      endDate.Format("2006-01-02"),
			"Preferences": pref,
		},
	}

	err = h.templates.Render(w, "statistics.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 統計データを取得
func (h *StatisticsHandler) GetStatisticsData(w http.ResponseWriter, r *http.Request) {
	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// クエリパラメータから期間を取得
	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")

	// 期間のバリデーション
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		http.Error(w, "無効な開始日", http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		http.Error(w, "無効な終了日", http.StatusBadRequest)
		return
	}

	// 統計データの取得
	stats, err := h.service.Record().GetStatistics(r.Context(), userID, start, end)
	if err != nil {
		http.Error(w, "統計データの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// JSONレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// 週間統計を取得
func (h *StatisticsHandler) GetWeeklyStats(w http.ResponseWriter, r *http.Request) {
	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// 直近の週間統計を取得
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)

	stats, err := h.service.Record().GetWeeklyStats(r.Context(), userID, startDate, endDate)
	if err != nil {
		http.Error(w, "週間統計の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// JSONレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// 月間統計を取得
func (h *StatisticsHandler) GetMonthlyStats(w http.ResponseWriter, r *http.Request) {
	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// 直近の月間統計を取得
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0)

	stats, err := h.service.Record().GetMonthlyStats(r.Context(), userID, startDate, endDate)
	if err != nil {
		http.Error(w, "月間統計の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// JSONレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

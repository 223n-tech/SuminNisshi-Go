package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// StatisticsHandler 統計情報画面のハンドラー
type StatisticsHandler struct {
	templates *TemplateManager
}

// StatisticsData 統計データの構造体
type StatisticsData struct {
	StartDate       time.Time              `json:"startDate"`
	EndDate         time.Time              `json:"endDate"`
	SleepTimes      []SleepTimeData        `json:"sleepTimes"`
	ScoreDistrib    []ScoreDistribution    `json:"scoreDistribution"`
	WeekdayAverages []WeekdayAverageData   `json:"weekdayAverages"`
	MonthlySummary  MonthlySummaryData     `json:"monthlySummary"`
}

// SleepTimeData 睡眠時間データ
type SleepTimeData struct {
	Date      time.Time `json:"date"`
	Duration  float64   `json:"duration"`    // 睡眠時間（時間）
	BedTime   string    `json:"bedTime"`     // 就寝時刻
	WakeTime  string    `json:"wakeTime"`    // 起床時刻
	Score     int       `json:"score"`       // 睡眠スコア
}

// ScoreDistribution スコア分布データ
type ScoreDistribution struct {
	Range string `json:"range"`  // スコア範囲（例: "60-69"）
	Count int    `json:"count"`  // 該当する記録の数
}

// WeekdayAverageData 曜日別平均データ
type WeekdayAverageData struct {
	Weekday       string  `json:"weekday"`      // 曜日
	AvgDuration   float64 `json:"avgDuration"`  // 平均睡眠時間
	AvgScore      float64 `json:"avgScore"`     // 平均スコア
}

// MonthlySummaryData 月間サマリーデータ
type MonthlySummaryData struct {
	AvgDuration    float64 `json:"avgDuration"`     // 平均睡眠時間
	AvgScore       float64 `json:"avgScore"`        // 平均スコア
	AvgBedTime     string  `json:"avgBedTime"`      // 平均就寝時刻
	AvgWakeTime    string  `json:"avgWakeTime"`     // 平均起床時刻
	DurationChange float64 `json:"durationChange"`   // 前月比（睡眠時間）
	ScoreChange    float64 `json:"scoreChange"`      // 前月比（スコア）
}

// NewStatisticsHandler 統計情報ハンドラーを作成
func NewStatisticsHandler(templates *TemplateManager) *StatisticsHandler {
	return &StatisticsHandler{
		templates: templates,
	}
}

// RegisterRoutes ルートの登録
func (h *StatisticsHandler) RegisterRoutes(r chi.Router) {
	r.Get("/statistics", h.Statistics)
	r.Get("/api/statistics/data", h.GetStatisticsData)
}

// Statistics 統計情報ページの表示
func (h *StatisticsHandler) Statistics(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:      "統計情報",
		ActiveMenu: "statistics",
	}

	err := h.templates.Render(w, "statistics.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetStatisticsData 統計データのJSONを返す
func (h *StatisticsHandler) GetStatisticsData(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータから期間を取得
	// TODO: 期間に基づくデータ取得の実装
	_ = r.URL.Query().Get("start")
	_ = r.URL.Query().Get("end")

	// TODO: 実際のデータベースからデータを取得
	// ここではサンプルデータを返す
	data := StatisticsData{
		SleepTimes: []SleepTimeData{
			{
				Date:     time.Now().AddDate(0, 0, -1),
				Duration: 7.5,
				BedTime:  "23:00",
				WakeTime: "06:30",
				Score:    85,
			},
			// 他のデータ...
		},
		ScoreDistrib: []ScoreDistribution{
			{Range: "90-100", Count: 3},
			{Range: "80-89", Count: 12},
			{Range: "70-79", Count: 8},
			{Range: "60-69", Count: 5},
			{Range: "0-59", Count: 2},
		},
		WeekdayAverages: []WeekdayAverageData{
			{Weekday: "月", AvgDuration: 7.2, AvgScore: 82},
			{Weekday: "火", AvgDuration: 7.0, AvgScore: 80},
			// 他の曜日...
		},
		MonthlySummary: MonthlySummaryData{
			AvgDuration:    7.2,
			AvgScore:       85,
			AvgBedTime:     "23:30",
			AvgWakeTime:    "06:45",
			DurationChange: 0.3,
			ScoreChange:    2.0,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// SleepRecordHandler 睡眠記録画面のハンドラー
type SleepRecordHandler struct {
	templates *TemplateManager
}

// SleepRecord 睡眠記録のデータ構造
type SleepRecord struct {
	ID             int64
	Date           time.Time
	BedTime        string
	WakeTime       string
	Duration       float64
	Score          int
	Quality        string
	ScoreColorClass string
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// SleepRecordFilter 睡眠記録のフィルター条件
type SleepRecordFilter struct {
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	DurationFrom float64   `json:"durationFrom"`
	DurationTo   float64   `json:"durationTo"`
	ScoreFrom    int       `json:"scoreFrom"`
	ScoreTo      int       `json:"scoreTo"`
	SortBy       string    `json:"sortBy"`
	SortOrder    string    `json:"sortOrder"`
}

// NewSleepRecordHandler 睡眠記録ハンドラーを作成
func NewSleepRecordHandler(templates *TemplateManager) *SleepRecordHandler {
	return &SleepRecordHandler{
		templates: templates,
	}
}

// RegisterRoutes ルートの登録
func (h *SleepRecordHandler) RegisterRoutes(r chi.Router) {
	r.Get("/sleep-records", h.List)
	r.Get("/sleep-records/new", h.New)
	r.Post("/sleep-records", h.Create)
	r.Get("/sleep-records/{id}", h.Show)
	r.Get("/sleep-records/{id}/edit", h.Edit)
	r.Put("/sleep-records/{id}", h.Update)
	r.Delete("/sleep-records/{id}", h.Delete)
	
	// API endpoints
	r.Get("/api/sleep-records", h.ListAPI)
	r.Post("/api/sleep-records/filter", h.FilterAPI)
}

// List 睡眠記録一覧の表示
func (h *SleepRecordHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: データベースから実際のデータを取得
	records := []SleepRecord{
		{
			ID:              1,
			Date:            time.Now().AddDate(0, 0, -1),
			BedTime:         "23:00",
			WakeTime:        "6:30",
			Duration:        7.5,
			Score:           85,
			Quality:         "良好",
			ScoreColorClass: "bg-success",
			Notes:           "よく眠れた",
		},
		{
			ID:              2,
			Date:            time.Now().AddDate(0, 0, -2),
			BedTime:         "23:30",
			WakeTime:        "6:45",
			Duration:        7.25,
			Score:           75,
			Quality:         "普通",
			ScoreColorClass: "bg-warning",
			Notes:           "途中で目が覚めた",
		},
	}

	data := &TemplateData{
		Title:      "睡眠記録一覧",
		ActiveMenu: "sleep-records",
		Data: map[string]interface{}{
			"Records": records,
		},
	}

	err := h.templates.Render(w, "sleep-records.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// New 新規記録フォームの表示
func (h *SleepRecordHandler) New(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:      "睡眠記録の作成",
		ActiveMenu: "sleep-records",
	}

	err := h.templates.Render(w, "sleep-records-form.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create 新規記録の作成
func (h *SleepRecordHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// TODO: バリデーションと保存処理の実装
	http.Redirect(w, r, "/sleep-records", http.StatusSeeOther)
}

// Show 記録詳細の表示
func (h *SleepRecordHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	recordID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	// TODO: データベースから記録を取得
	record := &SleepRecord{
		ID:              recordID,
		Date:            time.Now(),
		BedTime:         "23:00",
		WakeTime:        "6:30",
		Duration:        7.5,
		Score:           85,
		Quality:         "良好",
		ScoreColorClass: "bg-success",
		Notes:           "よく眠れた",
	}

	data := &TemplateData{
		Title:      "睡眠記録の詳細",
		ActiveMenu: "sleep-records",
		Data: map[string]interface{}{
			"Record": record,
		},
	}

	err = h.templates.Render(w, "sleep-records-detail.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Edit 記録編集フォームの表示
func (h *SleepRecordHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	recordID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	// TODO: データベースから記録を取得
	record := &SleepRecord{
		ID:       recordID,
		BedTime:  "23:00",
		WakeTime: "6:30",
	}

	data := &TemplateData{
		Title:      "睡眠記録の編集",
		ActiveMenu: "sleep-records",
		Data: map[string]interface{}{
			"Record": record,
		},
	}

	err = h.templates.Render(w, "sleep-records-form.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Update 記録の更新
func (h *SleepRecordHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// TODO: バリデーションと更新処理の実装
	http.Redirect(w, r, "/sleep-records", http.StatusSeeOther)
}

// Delete 記録の削除
func (h *SleepRecordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	// TODO: 削除処理の実装
	w.WriteHeader(http.StatusNoContent)
}

// ListAPI 睡眠記録一覧のJSON返却
func (h *SleepRecordHandler) ListAPI(w http.ResponseWriter, r *http.Request) {
	// TODO: データベースからデータを取得
	records := []SleepRecord{
		{
			ID:       1,
			BedTime:  "23:00",
			WakeTime: "6:30",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// FilterAPI フィルター条件に基づく睡眠記録の取得
func (h *SleepRecordHandler) FilterAPI(w http.ResponseWriter, r *http.Request) {
	var filter SleepRecordFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, "無効なフィルター条件", http.StatusBadRequest)
		return
	}

	// TODO: フィルター条件に基づくデータ取得の実装
	records := []SleepRecord{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// getScoreColorClass スコアに応じた色クラスを返す
func getScoreColorClass(score int) string {
	switch {
	case score >= 90:
		return "bg-success"
	case score >= 70:
		return "bg-info"
	case score >= 50:
		return "bg-warning"
	default:
		return "bg-danger"
	}
}

// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/sleep_records.go
// sleep_recordsは、睡眠記録関連のハンドラーを提供します。

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
	"github.com/223n-tech/SuiminNisshi-Go/internal/util"
	"github.com/go-chi/chi/v5"
)

// 睡眠記録関連のハンドラー
type SleepRecordHandler struct {
	templates *TemplateManager
	service   *service.Service
}

// SleepRecordHandlerを作成
func NewSleepRecordHandler(templates *TemplateManager, svc *service.Service) *SleepRecordHandler {
	return &SleepRecordHandler{
		templates: templates,
		service:   svc,
	}
}

// ルーティングを登録
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

// 睡眠記録一覧の表示
func (h *SleepRecordHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	// 睡眠状態のマスターデータを取得
	states, err := h.service.Record().GetStatesList(r.Context())
	if err != nil {
		http.Error(w, "睡眠状態の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// 食事種別のマスターデータを取得
	mealTypes, err := h.service.Record().GetMealTypesList(r.Context())
	if err != nil {
		http.Error(w, "食事種別の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// デフォルトの期間を設定（直近7日間）
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)

	// 睡眠記録の取得
	records, err := h.service.Record().GetRecordsByDateRange(r.Context(), userID, startDate, endDate)
	if err != nil {
		http.Error(w, "睡眠記録の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	data := &TemplateData{
		Title:      "睡眠記録一覧",
		ActiveMenu: "sleep-records",
		Data: map[string]interface{}{
			"Records":    records,
			"States":     states,
			"MealTypes":  mealTypes,
			"StartDate":  startDate.Format("2006-01-02"),
			"EndDate":    endDate.Format("2006-01-02"),
		},
	}

	err = h.templates.Render(w, "sleep-records.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 新規記録画面を表示
func (h *SleepRecordHandler) New(w http.ResponseWriter, r *http.Request) {
	states, err := h.service.Record().GetStatesList(r.Context())
	if err != nil {
		http.Error(w, "睡眠状態の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	mealTypes, err := h.service.Record().GetMealTypesList(r.Context())
	if err != nil {
		http.Error(w, "食事種別の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	data := &TemplateData{
		Title:      "睡眠記録の作成",
		ActiveMenu: "sleep-records",
		Data: map[string]interface{}{
			"States":     states,
			"MealTypes":  mealTypes,
			"Record": &models.SleepRecord{
				RecordDate: time.Now(),
			},
		},
	}

	err = h.templates.Render(w, "sleep-records-form.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 新規記録の作成
func (h *SleepRecordHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// フォームデータの取得と変換
	record := &models.SleepRecord{
		RecordDate:   util.ParseDate(r.FormValue("record_date")),
		TimeSlot:     util.ParseTime(r.FormValue("time_slot")),
		SleepStateID: util.ParseInt64(r.FormValue("sleep_state_id")),
		RecordType:   r.FormValue("record_type"),
		Note:         sql.NullString{String: r.FormValue("note"), Valid: true},
	}

	// 食事種別IDの設定（設定されている場合）
	if mealTypeID := r.FormValue("meal_type_id"); mealTypeID != "" {
		record.MealTypeID = sql.NullInt64{
			Int64: util.ParseInt64(mealTypeID),
			Valid: true,
		}
	}

	// 記録の作成
	err := h.service.Record().CreateRecord(r.Context(), record)
	if err != nil {
		http.Error(w, "睡眠記録の作成に失敗しました", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/sleep-records", http.StatusSeeOther)
}

// 記録の詳細表示
func (h *SleepRecordHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := util.ParseInt64(chi.URLParam(r, "id"))
	if id == 0 {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	record, err := h.service.Record().GetRecordWithRelations(r.Context(), id)
	if err != nil {
		http.Error(w, "睡眠記録の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	if record == nil {
		http.Error(w, "記録が見つかりません", http.StatusNotFound)
		return
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

// 記録の編集画面を表示
func (h *SleepRecordHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := util.ParseInt64(chi.URLParam(r, "id"))
	if id == 0 {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	record, err := h.service.Record().GetRecordWithRelations(r.Context(), id)
	if err != nil {
		http.Error(w, "睡眠記録の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	states, err := h.service.Record().GetStatesList(r.Context())
	if err != nil {
		http.Error(w, "睡眠状態の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	mealTypes, err := h.service.Record().GetMealTypesList(r.Context())
	if err != nil {
		http.Error(w, "食事種別の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	data := &TemplateData{
		Title:      "睡眠記録の編集",
		ActiveMenu: "sleep-records",
		Data: map[string]interface{}{
			"Record":     record,
			"States":     states,
			"MealTypes":  mealTypes,
		},
	}

	err = h.templates.Render(w, "sleep-records-form.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 記録の更新
func (h *SleepRecordHandler) Update(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
        return
    }

    id := util.ParseInt64(chi.URLParam(r, "id"))
    if id == 0 {
        http.Error(w, "無効なID", http.StatusBadRequest)
        return
    }

    // 既存の記録を取得
	/*
    recordWithRelations, err := h.service.Record().GetRecordWithRelations(r.Context(), id)
    if err != nil {
        http.Error(w, "睡眠記録の取得に失敗しました", http.StatusInternalServerError)
        return
    }
	*/

    // SleepRecord型の新しい変数を作成
    updatedRecord := &models.SleepRecord{
        ID:           id,
        RecordDate:   util.ParseDate(r.FormValue("record_date")),
        TimeSlot:     util.ParseTime(r.FormValue("time_slot")),
        SleepStateID: util.ParseInt64(r.FormValue("sleep_state_id")),
        RecordType:   r.FormValue("record_type"),
        Note:         sql.NullString{String: r.FormValue("note"), Valid: true},
    }

    if mealTypeID := r.FormValue("meal_type_id"); mealTypeID != "" {
        updatedRecord.MealTypeID = sql.NullInt64{
            Int64: util.ParseInt64(mealTypeID),
            Valid: true,
        }
    }

    // 記録の更新
    err := h.service.Record().UpdateRecord(r.Context(), updatedRecord)
    if err != nil {
        http.Error(w, "睡眠記録の更新に失敗しました", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/sleep-records", http.StatusSeeOther)
}

// 記録の削除
func (h *SleepRecordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := util.ParseInt64(chi.URLParam(r, "id"))
	if id == 0 {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	err := h.service.Record().DeleteRecord(r.Context(), id)
	if err != nil {
		http.Error(w, "睡眠記録の削除に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// 睡眠記録一覧のAPI
func (h *SleepRecordHandler) ListAPI(w http.ResponseWriter, r *http.Request) {
	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	records, err := h.service.Record().GetRecordsByDateRange(r.Context(), userID, startDate, endDate)
	if err != nil {
		http.Error(w, "睡眠記録の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// 睡眠記録のフィルターリングAPI
func (h *SleepRecordHandler) FilterAPI(w http.ResponseWriter, r *http.Request) {
	var filter service.SleepRecordFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, "無効なフィルター条件", http.StatusBadRequest)
		return
	}

	// TODO: 実際のユーザーIDを使用
	var userID int64 = 1 // 開発用

	records, err := h.service.Record().FilterRecords(r.Context(), userID, filter)
	if err != nil {
		http.Error(w, "睡眠記録のフィルターリングに失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// 睡眠スコアに応じた色のクラスを取得
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

// 時間範囲の妥当性をチェック
func validateTimeRange(startTime, endTime time.Time) error {
	if startTime.After(endTime) {
		return service.ErrInvalidTimeRange
	}

	duration := endTime.Sub(startTime)
	if duration > 24*time.Hour {
		return service.ErrTimeRangeTooLong
	}

	return nil
}


// 睡眠記録のバリデーション
func validateSleepRecord(record *models.SleepRecord) error {
	if record.RecordDate.IsZero() {
		return service.ErrInvalidDate
	}

	if record.TimeSlot.IsZero() {
		return service.ErrInvalidTimeSlot
	}

	if record.SleepStateID == 0 {
		return service.ErrInvalidSleepState
	}

	if record.RecordType == "" {
		return service.ErrInvalidRecordType
	}

	return nil
}

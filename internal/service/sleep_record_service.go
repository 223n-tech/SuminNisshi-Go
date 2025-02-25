// internal/service/sleep_record_service.go
// sleep_record_serviceは、睡眠記録関連のサービスを提供します。

// Package service provides application services.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

var (
	// ErrInvalidTimeRange 開始時刻は、終了時刻より前でなければなりません。
    ErrInvalidTimeRange = errors.New("start time must be before end time / 開始時刻は、終了時刻より前でなければなりません")
	// ErrTimeRangeTooLong 時間幅は24時間を超えてはいけません
    ErrTimeRangeTooLong = errors.New("time range cannot exceed 24 hours / 時間幅は24時間を超えてはいけません")
	// ErrInvalidDate 無効な日付です
    ErrInvalidDate = errors.New("invalid date / 無効な日付です")
	// ErrInvalidTimeSlot 無効なタイムスロットです
    ErrInvalidTimeSlot = errors.New("invalid time slot / 無効なタイムスロットです")
	// ErrInvalidSleepState 無効な睡眠状態です
    ErrInvalidSleepState = errors.New("invalid sleep state / 無効な睡眠状態です")
	// invalid record type 無効なレコード種別です
    ErrInvalidRecordType = errors.New("invalid record type / 無効なレコード種別です")
)

// 睡眠記録関連のサービス
type SleepRecordService struct {
	s *Service
}

// ダッシュボードに表示するデータ
type DashboardStats struct {
    TotalSleepHours   float64
    AverageSleepHours float64
    SleepQualityScore int
    TargetAchievement int
    RecentRecords     []models.SleepRecord
}

// フィルター条件
type SleepRecordFilter struct {
    StartDate string
    EndDate   string
    StateID   int64
    // その他のフィルター条件
}

// 新しいSleepRecordServiceを作成
func NewSleepRecordService(s *Service) *SleepRecordService {
	return &SleepRecordService{s: s}
}

// 新規睡眠記録を作成
func (s *SleepRecordService) CreateRecord(ctx context.Context, record *models.SleepRecord) error {
	// 時間枠の妥当性チェック
	if !record.IsValidTimeSlot() {
		return errors.New("invalid time slot")
	}

	// 睡眠状態の存在チェック
	state, err := s.s.repo.SleepState().GetByID(ctx, record.SleepStateID)
	if err != nil {
		return err
	}
	if state == nil {
		return errors.New("invalid sleep state")
	}

	// 食事種別の存在チェック（設定されている場合）
	if record.MealTypeID.Valid {
		mealType, err := s.s.repo.MealType().GetByID(ctx, record.MealTypeID.Int64)
		if err != nil {
			return err
		}
		if mealType == nil {
			return errors.New("invalid meal type")
		}
	}

	return s.s.repo.SleepRecord().Create(ctx, record)
}

// 日誌の全睡眠記録を取得
func (s *SleepRecordService) GetDiaryRecords(ctx context.Context, diaryID int64) ([]*models.SleepRecord, error) {
	return s.s.repo.SleepRecord().GetByDiaryID(ctx, diaryID)
}

// 関連データを含む睡眠記録を取得
func (s *SleepRecordService) GetRecordWithRelations(ctx context.Context, recordID int64) (*models.SleepRecordWithRelations, error) {
	return s.s.repo.SleepRecord().GetWithRelations(ctx, recordID)
}

// 日付範囲で睡眠記録を取得
func (s *SleepRecordService) GetRecordsByDateRange(ctx context.Context, diaryID int64, startDate, endDate time.Time) ([]*models.SleepRecord, error) {
	return s.s.repo.SleepRecord().GetByDateRange(ctx, diaryID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
}

// 睡眠記録を更新
func (s *SleepRecordService) UpdateRecord(ctx context.Context, record *models.SleepRecord) error {
	if !record.IsValidTimeSlot() {
		return errors.New("invalid time slot")
	}

	existing, err := s.s.repo.SleepRecord().GetByID(ctx, record.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("record not found")
	}

	return s.s.repo.SleepRecord().Update(ctx, record)
}

// 睡眠記録を削除
func (s *SleepRecordService) DeleteRecord(ctx context.Context, recordID int64) error {
	existing, err := s.s.repo.SleepRecord().GetByID(ctx, recordID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("record not found")
	}

	return s.s.repo.SleepRecord().Delete(ctx, recordID)
}

// 複数の睡眠記録を一括作成
func (s *SleepRecordService) BulkCreateRecords(ctx context.Context, records []*models.SleepRecord) error {
	for _, record := range records {
		if !record.IsValidTimeSlot() {
			return errors.New("invalid time slot")
		}
	}

	return s.s.repo.SleepRecord().BulkCreate(ctx, records)
}

// すべての睡眠状態を取得
func (s *SleepRecordService) GetStatesList(ctx context.Context) ([]*models.SleepState, error) {
	return s.s.repo.SleepState().GetAll(ctx)
}

// すべての食事種別を取得
func (s *SleepRecordService) GetMealTypesList(ctx context.Context) ([]*models.MealType, error) {
	return s.s.repo.MealType().GetAll(ctx)
}

// 時間範囲の妥当性をチェック
func (s *SleepRecordService) ValidateTimeRange(startTime, endTime time.Time) error {
	if startTime.After(endTime) {
		return errors.New("start time must be before end time")
	}

	duration := endTime.Sub(startTime)
	if duration > 24*time.Hour {
		return errors.New("time range cannot exceed 24 hours")
	}

	return nil
}

// ダッシュボードステータスを取得
func (s *SleepRecordService) GetDashboardStats(ctx context.Context, userID int64) (*DashboardStats, error) {
    // TODO: 実装
    return &DashboardStats{
        TotalSleepHours:   0,
        AverageSleepHours: 0,
        SleepQualityScore: 0,
        TargetAchievement: 0,
        RecentRecords:     []models.SleepRecord{},
    }, nil
}

// JSONデータに変換して取得
func (d *DashboardStats) ToJSON() []byte {
    jsonData, _ := json.Marshal(d)
    return jsonData
}

// 統計データを取得
func (s *SleepRecordService) GetStatistics(ctx context.Context, userID int64, startDate, endDate time.Time) (interface{}, error) {
    // 実装
    return nil, nil
}

// 週間データを取得
func (s *SleepRecordService) GetWeeklyStats(ctx context.Context, userID int64, startDate, endDate time.Time) (interface{}, error) {
    // 実装
    return nil, nil
}

// 月間データを取得
func (s *SleepRecordService) GetMonthlyStats(ctx context.Context, userID int64, startDate, endDate time.Time) (interface{}, error) {
    // 実装
    return nil, nil
}

// すべてのデータを取得
func (s *SleepRecordService) GetAllRecords(ctx context.Context, userID int64) ([]*models.SleepRecord, error) {
    // 実装
    return nil, nil
}

// 絞り込み条件で検索したデータを取得
func (s *SleepRecordService) FilterRecords(ctx context.Context, userID int64, filter SleepRecordFilter) ([]*models.SleepRecord, error) {
    // 実装
    return nil, nil
}

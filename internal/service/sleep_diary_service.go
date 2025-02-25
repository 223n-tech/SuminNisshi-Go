// internal/service/sleep_diary_service.go
// sleep_diary_serviceは、睡眠日誌関連のサービスを提供します

// Package service provides application services.
package service

import (
	"context"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

// 睡眠日誌関連のサービス
type SleepDiaryService struct {
	s *Service
}

// 新しいSleepDiaryServiceを作成
func NewSleepDiaryService(s *Service) *SleepDiaryService {
	return &SleepDiaryService{s: s}
}

// 新規睡眠日誌を作成
func (s *SleepDiaryService) CreateDiary(ctx context.Context, userID int64, startDate, endDate time.Time, name string) (*models.SleepDiary, error) {
	diary := &models.SleepDiary{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		DiaryName: name,
	}

	if err := s.s.repo.SleepDiary().Create(ctx, diary); err != nil {
		return nil, err
	}

	return diary, nil
}

// ユーザーの全睡眠日誌を取得
func (s *SleepDiaryService) GetUserDiaries(ctx context.Context, userID int64) ([]*models.SleepDiary, error) {
	return s.s.repo.SleepDiary().GetByUserID(ctx, userID)
}

// 日付範囲で睡眠日誌を取得
func (s *SleepDiaryService) GetDiaryByDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*models.SleepDiary, error) {
	return s.s.repo.SleepDiary().GetByDateRange(ctx, userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
}

// 睡眠日誌を更新
func (s *SleepDiaryService) UpdateDiary(ctx context.Context, diary *models.SleepDiary) error {
	return s.s.repo.SleepDiary().Update(ctx, diary)
}

// 睡眠日誌を削除
func (s *SleepDiaryService) DeleteDiary(ctx context.Context, diaryID int64) error {
	// 関連する睡眠記録も含めてトランザクションで削除
	return s.s.Transaction(ctx, func(ctx context.Context) error {
		records, err := s.s.repo.SleepRecord().GetByDiaryID(ctx, diaryID)
		if err != nil {
			return err
		}

		for _, record := range records {
			if err := s.s.repo.SleepRecord().Delete(ctx, record.ID); err != nil {
				return err
			}
		}

		return s.s.repo.SleepDiary().Delete(ctx, diaryID)
	})
}

// 睡眠日誌のサマリー情報を取得
func (s *SleepDiaryService) GetDiarySummary(ctx context.Context, diaryID int64) (*models.PDFStatistics, error) {
	_, err := s.s.repo.SleepRecord().GetByDiaryID(ctx, diaryID)
	if err != nil {
		return nil, err
	}

	diary, err := s.s.repo.SleepDiary().GetByID(ctx, diaryID)
	if err != nil {
		return nil, err
	}

	// 統計データの計算
	stats := &models.PDFStatistics{
		StartDate: diary.StartDate,
		EndDate:   diary.EndDate,
		TotalDays: diary.CalculateDuration(),
	}

	// TODO: 睡眠記録から統計値を計算
	// stats.AverageDuration = ...
	// stats.AverageBedTime = ...
	// stats.AverageWakeTime = ...
	// stats.AverageScore = ...

	return stats, nil
}

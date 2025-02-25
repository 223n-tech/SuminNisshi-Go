// internal/service/pdf_service.go
// pdf_serviceは、PDF出力関連のサービスを提供します。

// Package service provides application services.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

// PDF出力関連のサービス
type PDFService struct {
	s *Service
}

// 新しいPDFServiceを作成
func NewPDFService(s *Service) *PDFService {
	return &PDFService{s: s}
}

// 睡眠日誌のPDFを生成
func (s *PDFService) GenerateSleepDiaryPDF(ctx context.Context, userID, diaryID int64) ([]byte, error) {
	// 日誌の存在確認
	diary, err := s.s.repo.SleepDiary().GetByID(ctx, diaryID)
	if err != nil {
		return nil, err
	}
	if diary == nil {
		return nil, errors.New("diary not found")
	}

	// ユーザーの確認
	if diary.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	// ユーザー情報の取得
	user, err := s.s.repo.User().GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 睡眠記録の取得
	records, err := s.s.repo.SleepRecord().GetByDiaryID(ctx, diaryID)
	if err != nil {
		return nil, err
	}

	// 睡眠状態の取得
	states, err := s.s.repo.SleepState().GetAll(ctx)
	if err != nil {
		return nil, err
	}
	statesMap := make(map[int64]models.SleepState)
	for _, state := range states {
		statesMap[state.ID] = *state
	}

	// 食事種別の取得
	mealTypes, err := s.s.repo.MealType().GetAll(ctx)
	if err != nil {
		return nil, err
	}
	mealTypesMap := make(map[int64]models.MealType)
	for _, mealType := range mealTypes {
		mealTypesMap[mealType.ID] = *mealType
	}

	// ユーザーの睡眠設定を取得
	pref, err := s.s.repo.UserSleepPreference().GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// PDF出力用データの作成
	data := &models.PDFExportData{
		User:        *user,
		SleepDiary:  *diary,
		Records:     records,
		States:      statesMap,
		MealTypes:   mealTypesMap,
		Preferences: *pref,
	}

	// PDF出力用テンプレートの設定
	template := models.PDFTemplate{
		FontPath:   "internal/assets/fonts/ipaexg.ttf",
		PageWidth:  595.28,
		PageHeight: 841.89,
		Margin:     20,
	}

	// PDFの生成
	return data.GeneratePDF(template)
}

// 統計情報のPDFを生成
func (s *PDFService) GenerateStatisticsPDF(ctx context.Context, userID int64, startDate, endDate time.Time) ([]byte, error) {
	// 期間の妥当性チェック
	if startDate.After(endDate) {
		return nil, errors.New("invalid date range")
	}

	// 日誌の取得
	diaries, err := s.s.repo.SleepDiary().GetByDateRange(ctx, userID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	// ユーザー情報の取得
	user, err := s.s.repo.User().GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 睡眠記録の取得とデータの集計
	var allRecords []*models.SleepRecord
	for _, diary := range diaries {
		records, err := s.s.repo.SleepRecord().GetByDiaryID(ctx, diary.ID)
		if err != nil {
			return nil, err
		}
		allRecords = append(allRecords, records...)
	}

	// マスターデータの取得
	states, err := s.s.repo.SleepState().GetAll(ctx)
	if err != nil {
		return nil, err
	}
	statesMap := make(map[int64]models.SleepState)
	for _, state := range states {
		statesMap[state.ID] = *state
	}

	mealTypes, err := s.s.repo.MealType().GetAll(ctx)
	if err != nil {
		return nil, err
	}
	mealTypesMap := make(map[int64]models.MealType)
	for _, mealType := range mealTypes {
		mealTypesMap[mealType.ID] = *mealType
	}

	// ユーザーの睡眠設定を取得
	pref, err := s.s.repo.UserSleepPreference().GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 統計データの計算
	// stats := s.calculateStatistics(allRecords, startDate, endDate)

	// PDF出力用データの作成
	data := &models.PDFExportData{
		User: *user,
		SleepDiary: models.SleepDiary{
			StartDate: startDate,
			EndDate:   endDate,
		},
		Records:     allRecords,
		States:      statesMap,
		MealTypes:   mealTypesMap,
		Preferences: *pref,
	}

	// PDF出力用テンプレートの設定
	template := models.PDFTemplate{
		FontPath:   "internal/assets/fonts/ipaexg.ttf",
		PageWidth:  595.28,
		PageHeight: 841.89,
		Margin:     20,
	}

	// PDFの生成
	return data.GeneratePDF(template)
}

// 睡眠記録から統計データを計算
func (s *PDFService) calculateStatistics(_ []*models.SleepRecord, startDate, endDate time.Time) *models.PDFStatistics {
	stats := &models.PDFStatistics{
		StartDate: startDate,
		EndDate:   endDate,
		TotalDays: int(endDate.Sub(startDate).Hours() / 24),
	}

	// TODO: 睡眠記録から各種統計値を計算
	// - 平均睡眠時間
	// - 平均就寝時刻
	// - 平均起床時刻
	// - 平均睡眠スコア
	// など

	return stats
}

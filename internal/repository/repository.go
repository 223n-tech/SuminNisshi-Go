// internal/repository/repository.go
// このパッケージは、データベースの操作を行うためのインターフェイスを提供します。
// これにより、データベースの実装を変更する際に、コードの変更を最小限に抑えることができます。

// Package repository provides interfaces for abstracting database operations.
package repository

import (
	"context"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

// 全リポジトリのインターフェイス
type Repository interface {
	// サブリポジトリの取得
	User() UserRepository
	SleepDiary() SleepDiaryRepository
	SleepRecord() SleepRecordRepository
	SleepState() SleepStateRepository
	MealType() MealTypeRepository
	UserSleepPreference() UserSleepPreferenceRepository
	// トランザクション
	Transaction(ctx context.Context, fn func(Repository) error) error
}

// ユーザー情報のリポジトリーインターフェイス
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
	UpdateLastLogin(ctx context.Context, id int64) error
}

// 睡眠日誌のリポジトリーインターフェイス
type SleepDiaryRepository interface {
	GetByID(ctx context.Context, id int64) (*models.SleepDiary, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.SleepDiary, error)
	GetByDateRange(ctx context.Context, userID int64, startDate, endDate string) ([]*models.SleepDiary, error)
	Create(ctx context.Context, diary *models.SleepDiary) error
	Update(ctx context.Context, diary *models.SleepDiary) error
	Delete(ctx context.Context, id int64) error
}

// 睡眠記録のリポジトリーインターフェイス
type SleepRecordRepository interface {
	GetByID(ctx context.Context, id int64) (*models.SleepRecord, error)
	GetByDiaryID(ctx context.Context, diaryID int64) ([]*models.SleepRecord, error)
	GetByDateRange(ctx context.Context, diaryID int64, startDate, endDate string) ([]*models.SleepRecord, error)
	GetWithRelations(ctx context.Context, id int64) (*models.SleepRecordWithRelations, error)
	Create(ctx context.Context, record *models.SleepRecord) error
	Update(ctx context.Context, record *models.SleepRecord) error
	Delete(ctx context.Context, id int64) error
	BulkCreate(ctx context.Context, records []*models.SleepRecord) error
}

// 睡眠状態のリポジトリーインターフェイス
type SleepStateRepository interface {
	GetByID(ctx context.Context, id int64) (*models.SleepState, error)
	GetAll(ctx context.Context) ([]*models.SleepState, error)
	GetByCode(ctx context.Context, code string) (*models.SleepState, error)
	Create(ctx context.Context, state *models.SleepState) error
	Update(ctx context.Context, state *models.SleepState) error
	Delete(ctx context.Context, id int64) error
}

// 食事種別のリポジトリーインターフェイス
type MealTypeRepository interface {
	GetByID(ctx context.Context, id int64) (*models.MealType, error)
	GetAll(ctx context.Context) ([]*models.MealType, error)
	GetByCode(ctx context.Context, code string) (*models.MealType, error)
	Create(ctx context.Context, mealType *models.MealType) error
	Update(ctx context.Context, mealType *models.MealType) error
	Delete(ctx context.Context, id int64) error
}

// ユーザー睡眠設定のリポジトリーインターフェイス
type UserSleepPreferenceRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*models.UserSleepPreference, error)
	Create(ctx context.Context, pref *models.UserSleepPreference) error
	Update(ctx context.Context, pref *models.UserSleepPreference) error
	Delete(ctx context.Context, userID int64) error
	GetDefaultPreference(userID int64) *models.UserSleepPreference
}

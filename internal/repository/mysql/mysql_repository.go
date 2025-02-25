// internal/repository/mysql/mysql_repository.go
// mysql_repositoryは、MySQLリポジトリの実装を提供します。

// Package mysql provides MySQL repository implementations.
package mysql

import (
	"context"
	"database/sql"

	"github.com/223n-tech/SuiminNisshi-Go/internal/repository"
)

// MySQLリポジトリの実装
type MySQLRepository struct {
	db *sql.DB
	tx *sql.Tx
}

// 新しいMySQLリポジトリを作成
func NewMySQLRepository(db *sql.DB) repository.Repository {
	return &MySQLRepository{
		db: db,
	}
}

// UserRepositoryを取得
func (r *MySQLRepository) User() repository.UserRepository {
	return &UserRepository{repo: r}
}

// SleepDiaryRepositoryを取得
func (r *MySQLRepository) SleepDiary() repository.SleepDiaryRepository {
	return &SleepDiaryRepository{repo: r}
}

// SleepRecordRepositoryを取得
func (r *MySQLRepository) SleepRecord() repository.SleepRecordRepository {
	return &SleepRecordRepository{repo: r}
}

// SleepStateRepositoryを取得
func (r *MySQLRepository) SleepState() repository.SleepStateRepository {
	return &SleepStateRepository{repo: r}
}

// MealTypeRepositoryを取得
func (r *MySQLRepository) MealType() repository.MealTypeRepository {
	return &MealTypeRepository{repo: r}
}

// UserSleepPreferenceRepositoryを取得
func (r *MySQLRepository) UserSleepPreference() repository.UserSleepPreferenceRepository {
	return &UserSleepPreferenceRepository{repo: r}
}

// トランザクションを実行
func (r *MySQLRepository) Transaction(ctx context.Context, fn func(repository.Repository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	repo := &MySQLRepository{
		db: r.db,
		tx: tx,
	}

	if err := fn(repo); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// アクティブなデータベース接続を取得
func (r *MySQLRepository) getDB() interface{} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

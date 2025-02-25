// internal/repository/mysql/sleep_diary_repository.go
// sleep_diary_repositoryは、睡眠日誌のリポジトリを提供します。

// Package mysql provides MySQL repository implementations.
package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

// SleepDiaryRepositoryのMySQL実装
type SleepDiaryRepository struct {
	repo *MySQLRepository
}

// IDで睡眠日誌を検索
func (r *SleepDiaryRepository) GetByID(ctx context.Context, id int64) (*models.SleepDiary, error) {
	query := `
		SELECT id, user_id, start_date, end_date, diary_name, note, created, modified, deleted
		FROM sleep_diaries
		WHERE id = ? AND deleted IS NULL
	`

	diary := &models.SleepDiary{}
	err := r.repo.getDB().(*sql.DB).QueryRowContext(ctx, query, id).Scan(
		&diary.ID,
		&diary.UserID,
		&diary.StartDate,
		&diary.EndDate,
		&diary.DiaryName,
		&diary.Note,
		&diary.Created,
		&diary.Modified,
		&diary.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return diary, nil
}

// ユーザーIDで睡眠日誌を検索
func (r *SleepDiaryRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.SleepDiary, error) {
	query := `
		SELECT id, user_id, start_date, end_date, diary_name, note, created, modified, deleted
		FROM sleep_diaries
		WHERE user_id = ? AND deleted IS NULL
		ORDER BY start_date DESC
	`

	rows, err := r.repo.getDB().(*sql.DB).QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diaries []*models.SleepDiary
	for rows.Next() {
		diary := &models.SleepDiary{}
		err := rows.Scan(
			&diary.ID,
			&diary.UserID,
			&diary.StartDate,
			&diary.EndDate,
			&diary.DiaryName,
			&diary.Note,
			&diary.Created,
			&diary.Modified,
			&diary.Deleted,
		)
		if err != nil {
			return nil, err
		}
		diaries = append(diaries, diary)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return diaries, nil
}

// 日付範囲で睡眠日誌を検索
func (r *SleepDiaryRepository) GetByDateRange(ctx context.Context, userID int64, startDate, endDate string) ([]*models.SleepDiary, error) {
	query := `
		SELECT id, user_id, start_date, end_date, diary_name, note, created, modified, deleted
		FROM sleep_diaries
		WHERE user_id = ? 
		AND deleted IS NULL
		AND (
			(start_date BETWEEN ? AND ?) 
			OR (end_date BETWEEN ? AND ?)
			OR (start_date <= ? AND end_date >= ?)
		)
		ORDER BY start_date
	`

	rows, err := r.repo.getDB().(*sql.DB).QueryContext(ctx, query,
		userID, startDate, endDate, startDate, endDate, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diaries []*models.SleepDiary
	for rows.Next() {
		diary := &models.SleepDiary{}
		err := rows.Scan(
			&diary.ID,
			&diary.UserID,
			&diary.StartDate,
			&diary.EndDate,
			&diary.DiaryName,
			&diary.Note,
			&diary.Created,
			&diary.Modified,
			&diary.Deleted,
		)
		if err != nil {
			return nil, err
		}
		diaries = append(diaries, diary)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return diaries, nil
}

// 新規睡眠日誌を作成
func (r *SleepDiaryRepository) Create(ctx context.Context, diary *models.SleepDiary) error {
	query := `
		INSERT INTO sleep_diaries (
			user_id, start_date, end_date, diary_name, note, created, modified
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		diary.UserID,
		diary.StartDate,
		diary.EndDate,
		diary.DiaryName,
		diary.Note,
		now,
		now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	diary.ID = id
	diary.Created = now
	diary.Modified = now

	return nil
}

// 睡眠日誌を更新
func (r *SleepDiaryRepository) Update(ctx context.Context, diary *models.SleepDiary) error {
	query := `
		UPDATE sleep_diaries
		SET start_date = ?, end_date = ?, diary_name = ?, note = ?, modified = ?
		WHERE id = ? AND deleted IS NULL
	`

	now := time.Now()
	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		diary.StartDate,
		diary.EndDate,
		diary.DiaryName,
		diary.Note,
		now,
		diary.ID,
	)

	if err != nil {
		return err
	}

	diary.Modified = now
	return nil
}

// 睡眠日誌を論理削除
func (r *SleepDiaryRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE sleep_diaries
		SET deleted = ?
		WHERE id = ? AND deleted IS NULL
	`

	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		time.Now(),
		id,
	)

	return err
}

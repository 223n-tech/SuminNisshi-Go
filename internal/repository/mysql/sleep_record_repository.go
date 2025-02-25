// internal/repository/mysql/sleep_record_repository.go
// sleep_record_repositoryは、睡眠記録のリポジトリを提供します。

// Package mysql provides MySQL repository implementations.
package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

// SleepRecordRepositoryのMySQL実装
type SleepRecordRepository struct {
	repo *MySQLRepository
}

// IDで睡眠記録を検索
func (r *SleepRecordRepository) GetByID(ctx context.Context, id int64) (*models.SleepRecord, error) {
	query := `
		SELECT id, sleep_diary_id, sleep_state_id, record_date, time_slot, record_type, meal_type_id, note, created, modified, deleted
		FROM sleep_records
		WHERE id = ? AND deleted IS NULL
	`

	record := &models.SleepRecord{}
	err := r.repo.getDB().(*sql.DB).QueryRowContext(ctx, query, id).Scan(
		&record.ID,
		&record.SleepDiaryID,
		&record.SleepStateID,
		&record.RecordDate,
		&record.TimeSlot,
		&record.RecordType,
		&record.MealTypeID,
		&record.Note,
		&record.Created,
		&record.Modified,
		&record.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return record, nil
}

// 日誌IDで睡眠記録を検索
func (r *SleepRecordRepository) GetByDiaryID(ctx context.Context, diaryID int64) ([]*models.SleepRecord, error) {
	query := `
		SELECT id, sleep_diary_id, sleep_state_id, record_date, time_slot, record_type, meal_type_id, note, created, modified, deleted
		FROM sleep_records
		WHERE sleep_diary_id = ? AND deleted IS NULL
		ORDER BY record_date, time_slot
	`

	rows, err := r.repo.getDB().(*sql.DB).QueryContext(ctx, query, diaryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.SleepRecord
	for rows.Next() {
		record := &models.SleepRecord{}
		err := rows.Scan(
			&record.ID,
			&record.SleepDiaryID,
			&record.SleepStateID,
			&record.RecordDate,
			&record.TimeSlot,
			&record.RecordType,
			&record.MealTypeID,
			&record.Note,
			&record.Created,
			&record.Modified,
			&record.Deleted,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// 日付範囲で睡眠記録を検索
func (r *SleepRecordRepository) GetByDateRange(ctx context.Context, diaryID int64, startDate, endDate string) ([]*models.SleepRecord, error) {
	query := `
		SELECT id, sleep_diary_id, sleep_state_id, record_date, time_slot, record_type, meal_type_id, note, created, modified, deleted
		FROM sleep_records
		WHERE sleep_diary_id = ?
		AND record_date BETWEEN ? AND ?
		AND deleted IS NULL
		ORDER BY record_date, time_slot
	`

	rows, err := r.repo.getDB().(*sql.DB).QueryContext(ctx, query, diaryID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.SleepRecord
	for rows.Next() {
		record := &models.SleepRecord{}
		err := rows.Scan(
			&record.ID,
			&record.SleepDiaryID,
			&record.SleepStateID,
			&record.RecordDate,
			&record.TimeSlot,
			&record.RecordType,
			&record.MealTypeID,
			&record.Note,
			&record.Created,
			&record.Modified,
			&record.Deleted,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// 関連データを含めて睡眠記録を検索
func (r *SleepRecordRepository) GetWithRelations(ctx context.Context, id int64) (*models.SleepRecordWithRelations, error) {
	query := `
		SELECT 
			r.id, r.sleep_diary_id, r.sleep_state_id, r.record_date, r.time_slot,
			r.record_type, r.meal_type_id, r.note, r.created, r.modified, r.deleted,
			s.id, s.state_name, s.state_code, s.state_description, s.display_symbol,
			s.display_order, s.created, s.modified, s.deleted,
			d.id, d.user_id, d.start_date, d.end_date, d.diary_name, d.note,
			d.created, d.modified, d.deleted,
			m.id, m.type_name, m.type_code, m.display_symbol, m.display_order,
			m.created, m.modified, m.deleted
		FROM sleep_records r
		LEFT JOIN sleep_states s ON r.sleep_state_id = s.id
		LEFT JOIN sleep_diaries d ON r.sleep_diary_id = d.id
		LEFT JOIN meal_types m ON r.meal_type_id = m.id
		WHERE r.id = ? AND r.deleted IS NULL
	`

	record := &models.SleepRecordWithRelations{}
	var mealType models.MealType
	err := r.repo.getDB().(*sql.DB).QueryRowContext(ctx, query, id).Scan(
		&record.ID, &record.SleepDiaryID, &record.SleepStateID, &record.RecordDate,
		&record.TimeSlot, &record.RecordType, &record.MealTypeID, &record.Note,
		&record.Created, &record.Modified, &record.Deleted,
		&record.State.ID, &record.State.StateName, &record.State.StateCode,
		&record.State.StateDescription, &record.State.DisplaySymbol,
		&record.State.DisplayOrder, &record.State.Created, &record.State.Modified,
		&record.State.Deleted,
		&record.Diary.ID, &record.Diary.UserID, &record.Diary.StartDate,
		&record.Diary.EndDate, &record.Diary.DiaryName, &record.Diary.Note,
		&record.Diary.Created, &record.Diary.Modified, &record.Diary.Deleted,
		&mealType.ID, &mealType.TypeName, &mealType.TypeCode, &mealType.DisplaySymbol,
		&mealType.DisplayOrder, &mealType.Created, &mealType.Modified, &mealType.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if record.MealTypeID.Valid {
		record.MealType = &mealType
	}

	return record, nil
}

// 新規睡眠記録を作成
func (r *SleepRecordRepository) Create(ctx context.Context, record *models.SleepRecord) error {
	query := `
		INSERT INTO sleep_records (
			sleep_diary_id, sleep_state_id, record_date, time_slot,
			record_type, meal_type_id, note, created, modified
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		record.SleepDiaryID,
		record.SleepStateID,
		record.RecordDate,
		record.TimeSlot,
		record.RecordType,
		record.MealTypeID,
		record.Note,
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

	record.ID = id
	record.Created = now
	record.Modified = now

	return nil
}

// 睡眠記録を更新
func (r *SleepRecordRepository) Update(ctx context.Context, record *models.SleepRecord) error {
	query := `
		UPDATE sleep_records
		SET sleep_state_id = ?, record_date = ?, time_slot = ?, record_type = ?, meal_type_id = ?, note = ?, modified = ?
		WHERE id = ? AND deleted IS NULL
	`

	now := time.Now()
	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		record.SleepStateID,
		record.RecordDate,
		record.TimeSlot,
		record.RecordType,
		record.MealTypeID,
		record.Note,
		now,
		record.ID,
	)

	if err != nil {
		return err
	}

	record.Modified = now
	return nil
}

// 睡眠記録を論理削除
func (r *SleepRecordRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE sleep_records
		SET deleted = ?
		WHERE id = ? AND deleted IS NULL
	`

	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		time.Now(),
		id,
	)

	return err
}

// 複数の睡眠記録を一括作成
func (r *SleepRecordRepository) BulkCreate(ctx context.Context, records []*models.SleepRecord) error {
	query := `
		INSERT INTO sleep_records (
			sleep_diary_id, sleep_state_id, record_date, time_slot,
			record_type, meal_type_id, note, created, modified
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	tx, err := r.repo.getDB().(*sql.DB).BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		result, err := stmt.ExecContext(ctx,
			record.SleepDiaryID,
			record.SleepStateID,
			record.RecordDate,
			record.TimeSlot,
			record.RecordType,
			record.MealTypeID,
			record.Note,
			now,
			now,
		)
		if err != nil {
			tx.Rollback()
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}

		record.ID = id
		record.Created = now
		record.Modified = now
	}

	return tx.Commit()
}

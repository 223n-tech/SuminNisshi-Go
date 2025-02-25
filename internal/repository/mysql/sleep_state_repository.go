// internal/repository/mysql/sleep_state_repository.go
// sleep_state_repositoryは、睡眠状態のリポジトリを提供します。

// Package mysql provides MySQL repository implementations.
package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

// SleepStateRepositoryのMySQL実装
type SleepStateRepository struct {
	repo *MySQLRepository
}

// IDで睡眠状態を検索
func (r *SleepStateRepository) GetByID(ctx context.Context, id int64) (*models.SleepState, error) {
	query := `
		SELECT id, state_name, state_code, state_description, display_symbol, display_order, created, modified, deleted
		FROM sleep_states
		WHERE id = ? AND deleted IS NULL
	`

	state := &models.SleepState{}
	err := r.repo.getDB().(*sql.DB).QueryRowContext(ctx, query, id).Scan(
		&state.ID,
		&state.StateName,
		&state.StateCode,
		&state.StateDescription,
		&state.DisplaySymbol,
		&state.DisplayOrder,
		&state.Created,
		&state.Modified,
		&state.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return state, nil
}

// すべての睡眠状態を検索
func (r *SleepStateRepository) GetAll(ctx context.Context) ([]*models.SleepState, error) {
	query := `
		SELECT id, state_name, state_code, state_description, display_symbol, display_order, created, modified, deleted
		FROM sleep_states
		WHERE deleted IS NULL
		ORDER BY display_order
	`

	rows, err := r.repo.getDB().(*sql.DB).QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*models.SleepState
	for rows.Next() {
		state := &models.SleepState{}
		err := rows.Scan(
			&state.ID,
			&state.StateName,
			&state.StateCode,
			&state.StateDescription,
			&state.DisplaySymbol,
			&state.DisplayOrder,
			&state.Created,
			&state.Modified,
			&state.Deleted,
		)
		if err != nil {
			return nil, err
		}
		states = append(states, state)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return states, nil
}

// 状態コードで睡眠状態を検索
func (r *SleepStateRepository) GetByCode(ctx context.Context, code string) (*models.SleepState, error) {
	query := `
		SELECT id, state_name, state_code, state_description, display_symbol, display_order, created, modified, deleted
		FROM sleep_states
		WHERE state_code = ? AND deleted IS NULL
	`

	state := &models.SleepState{}
	err := r.repo.getDB().(*sql.DB).QueryRowContext(ctx, query, code).Scan(
		&state.ID,
		&state.StateName,
		&state.StateCode,
		&state.StateDescription,
		&state.DisplaySymbol,
		&state.DisplayOrder,
		&state.Created,
		&state.Modified,
		&state.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return state, nil
}

// 新規睡眠状態を作成
func (r *SleepStateRepository) Create(ctx context.Context, state *models.SleepState) error {
	query := `
		INSERT INTO sleep_states (
			state_name, state_code, state_description, display_symbol,
			display_order, created, modified
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		state.StateName,
		state.StateCode,
		state.StateDescription,
		state.DisplaySymbol,
		state.DisplayOrder,
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

	state.ID = id
	state.Created = now
	state.Modified = now

	return nil
}

// 睡眠状態を更新
func (r *SleepStateRepository) Update(ctx context.Context, state *models.SleepState) error {
	query := `
		UPDATE sleep_states
		SET state_name = ?, state_code = ?, state_description = ?,
			display_symbol = ?, display_order = ?, modified = ?
		WHERE id = ? AND deleted IS NULL
	`

	now := time.Now()
	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		state.StateName,
		state.StateCode,
		state.StateDescription,
		state.DisplaySymbol,
		state.DisplayOrder,
		now,
		state.ID,
	)

	if err != nil {
		return err
	}

	state.Modified = now
	return nil
}

// 睡眠状態を論理削除
func (r *SleepStateRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE sleep_states
		SET deleted = ?
		WHERE id = ? AND deleted IS NULL
	`

	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		time.Now(),
		id,
	)

	return err
}

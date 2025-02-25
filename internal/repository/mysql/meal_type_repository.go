// internal/repository/mysql/meal_type_repository.go
// meal_type_repositoryは、食事種別のリポジトリを提供します。

// Package mysql provides MySQL repository implementations.
package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

// MealTypeRepositoryのMySQL実装
type MealTypeRepository struct {
	repo *MySQLRepository
}

// IDで食事種別を検索
func (r *MealTypeRepository) GetByID(ctx context.Context, id int64) (*models.MealType, error) {
	query := `
		SELECT id, type_name, type_code, display_symbol, display_order, created, modified, deleted
		FROM meal_types
		WHERE id = ? AND deleted IS NULL
	`

	mealType := &models.MealType{}
	err := r.repo.getDB().(*sql.DB).QueryRowContext(ctx, query, id).Scan(
		&mealType.ID,
		&mealType.TypeName,
		&mealType.TypeCode,
		&mealType.DisplaySymbol,
		&mealType.DisplayOrder,
		&mealType.Created,
		&mealType.Modified,
		&mealType.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mealType, nil
}

// すべての食事種別を検索
func (r *MealTypeRepository) GetAll(ctx context.Context) ([]*models.MealType, error) {
	query := `
		SELECT id, type_name, type_code, display_symbol, display_order, created, modified, deleted
		FROM meal_types
		WHERE deleted IS NULL
		ORDER BY display_order
	`

	rows, err := r.repo.getDB().(*sql.DB).QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mealTypes []*models.MealType
	for rows.Next() {
		mealType := &models.MealType{}
		err := rows.Scan(
			&mealType.ID,
			&mealType.TypeName,
			&mealType.TypeCode,
			&mealType.DisplaySymbol,
			&mealType.DisplayOrder,
			&mealType.Created,
			&mealType.Modified,
			&mealType.Deleted,
		)
		if err != nil {
			return nil, err
		}
		mealTypes = append(mealTypes, mealType)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mealTypes, nil
}

// 種別コードで食事種別を検索
func (r *MealTypeRepository) GetByCode(ctx context.Context, code string) (*models.MealType, error) {
	query := `
		SELECT id, type_name, type_code, display_symbol, display_order, created, modified, deleted
		FROM meal_types
		WHERE type_code = ? AND deleted IS NULL
	`

	mealType := &models.MealType{}
	err := r.repo.getDB().(*sql.DB).QueryRowContext(ctx, query, code).Scan(
		&mealType.ID,
		&mealType.TypeName,
		&mealType.TypeCode,
		&mealType.DisplaySymbol,
		&mealType.DisplayOrder,
		&mealType.Created,
		&mealType.Modified,
		&mealType.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mealType, nil
}

// 新規食事種別を作成
func (r *MealTypeRepository) Create(ctx context.Context, mealType *models.MealType) error {
	query := `
		INSERT INTO meal_types (
			type_name, type_code, display_symbol, display_order,
			created, modified
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		mealType.TypeName,
		mealType.TypeCode,
		mealType.DisplaySymbol,
		mealType.DisplayOrder,
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

	mealType.ID = id
	mealType.Created = now
	mealType.Modified = now

	return nil
}

// 食事種別を更新
func (r *MealTypeRepository) Update(ctx context.Context, mealType *models.MealType) error {
	query := `
		UPDATE meal_types
		SET type_name = ?, type_code = ?, display_symbol = ?, display_order = ?, modified = ?
		WHERE id = ? AND deleted IS NULL
	`

	now := time.Now()
	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		mealType.TypeName,
		mealType.TypeCode,
		mealType.DisplaySymbol,
		mealType.DisplayOrder,
		now,
		mealType.ID,
	)

	if err != nil {
		return err
	}

	mealType.Modified = now
	return nil
}

// 食事種別を論理削除
func (r *MealTypeRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE meal_types
		SET deleted = ?
		WHERE id = ? AND deleted IS NULL
	`

	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		time.Now(),
		id,
	)

	return err
}

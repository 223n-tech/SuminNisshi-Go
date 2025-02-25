// internal/repository/mysql/user_repository.go
// user_repositoryは、ユーザーのリポジトリを提供します。

// Package mysql provides MySQL repository implementations.
package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

// UserRepositoryのMySQL実装
type UserRepository struct {
	repo *MySQLRepository
}

// IDでユーザーを検索
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, email, display_name, password_hash, last_login_datetime, created, modified, deleted
		FROM users
		WHERE id = ? AND deleted IS NULL
	`

	user := &models.User{}
	err := r.repo.getDB().(*sql.DB).QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&user.LastLoginDatetime,
		&user.Created,
		&user.Modified,
		&user.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// メールアドレスでユーザーを検索
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, display_name, password_hash, last_login_datetime, created, modified, deleted
		FROM users
		WHERE email = ? AND deleted IS NULL
	`

	user := &models.User{}
	err := r.repo.getDB().(*sql.DB).QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&user.LastLoginDatetime,
		&user.Created,
		&user.Modified,
		&user.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// 新規ユーザーを作成
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			email, display_name, password_hash, created, modified
		) VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		user.Email,
		user.DisplayName,
		user.PasswordHash,
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

	user.ID = id
	user.Created = now
	user.Modified = now

	return nil
}

// ユーザー情報を更新
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = ?, display_name = ?, password_hash = ?, modified = ?
		WHERE id = ? AND deleted IS NULL
	`

	now := time.Now()
	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		user.Email,
		user.DisplayName,
		user.PasswordHash,
		now,
		user.ID,
	)

	if err != nil {
		return err
	}

	user.Modified = now
	return nil
}

// ユーザーを論理削除
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE users
		SET deleted = ?
		WHERE id = ? AND deleted IS NULL
	`

	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		time.Now(),
		id,
	)

	return err
}

// 最終ログイン日時を更新
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	query := `
		UPDATE users
		SET last_login_datetime = ?, modified = ?
		WHERE id = ? AND deleted IS NULL
	`

	now := time.Now()
	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		now,
		now,
		id,
	)

	return err
}

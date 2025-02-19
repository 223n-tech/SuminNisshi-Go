package repository

import (
	"database/sql"
	"fmt"

	"github.com/223n-tech/SuiminNisshi-Go/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

// RepositoryInterface はリポジトリの操作を定義するインターフェース
type RepositoryInterface interface {
	Close() error
	// TODO: 他のリポジトリメソッドを追加
}

// Repository はデータアクセス層の構造体
type Repository struct {
	db *sql.DB
}

// NewDB データベース接続を初期化
func NewDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// コネクションプールの設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}

// NewRepository リポジトリインスタンスを作成
func NewRepository(db *sql.DB) RepositoryInterface {
	return &Repository{
		db: db,
	}
}

// Close データベース接続を閉じる
func (r *Repository) Close() error {
	return r.db.Close()
}

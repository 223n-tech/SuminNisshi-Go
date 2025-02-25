// internal/repository/mysql/db.go
// dbは、データベース接続を提供します。

// Package mysql provides MySQL repository implementations.
package mysql

import (
	"database/sql"
	"fmt"
	"time"

	// MySQLドライバを使用するために必要
	_ "github.com/go-sql-driver/mysql"
)

// データベース接続設定
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// データベース接続を初期化
func NewDB(config DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// コネクションプールの設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

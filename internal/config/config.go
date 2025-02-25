// internal/config/config.go
// configは、アプリケーション全体の設定を保持する構造体を定義しています。

// Package config provides a structure to hold the application-wide configuration.
package config

import (
	"os"
	"strconv"
)

/*
	アプリケーション全体の設定を保持する構造体
*/
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

/*
	サーバー関連の設定
*/
type ServerConfig struct {
	Port    int
	Host    string
	BaseURL string
}

/*
	データベース関連の設定
*/
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

/*
	環境変数から設定を読み込む
*/
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:    getEnvInt("APP_PORT", 8080),
			Host:    getEnvStr("APP_HOST", "localhost"),
			BaseURL: getEnvStr("APP_BASE_URL", "http://localhost:8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnvStr("DB_HOST", "db"),
			Port:     getEnvInt("DB_PORT", 3306),
			User:     getEnvStr("DB_USER", "suiminnisshi"),
			Password: getEnvStr("DB_PASSWORD", "suiminnisshi_password"),
			DBName:   getEnvStr("DB_NAME", "suiminnisshi"),
		},
	}

	return cfg, nil
}

/*
	環境変数から文字列を取得
*/
func getEnvStr(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

/*
	環境変数から整数を取得
*/
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

package config

import (
	"os"
	"strconv"
)

// Config アプリケーション全体の設定を保持する構造体
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig サーバー関連の設定
type ServerConfig struct {
	Port    int
	Host    string
	BaseURL string
}

// DatabaseConfig データベース関連の設定
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// Load 環境変数から設定を読み込む
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

// getEnvStr 環境変数から文字列を取得
func getEnvStr(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt 環境変数から数値を取得
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

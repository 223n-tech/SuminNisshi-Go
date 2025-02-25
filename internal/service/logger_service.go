// internal/logger/custom_logger.go
// loggerパッケージは、ログ機能を拡張したパッケージです。

// package logger implements package that extends logging capabilities
package service

import (
	"log"
	"os"
)

// 標準ログパッケージを拡張したロガー
type LoggerService struct {
	*log.Logger
	level LogLevel
}

// ログのレベルを表す型
type LogLevel int

// ログレベルの定数
const (
	DebugLevel LogLevel = iota
	InfoLevel
	ErrorLevel
)

// 新しいLoggerServiceインスタンスを作成します
func NewLoggerService(level LogLevel, logger *log.Logger) *LoggerService {
	return &LoggerService{
		Logger: logger,
		level:  level,
	}
}

// デバッグレベルのログを出力します
func (l *LoggerService) Debug(format string, v ...interface{}) {
	if l.level <= DebugLevel {
		l.Logger.Printf("[DEBUG] "+format, v...)
	}
}

// 情報レベルのログを出力します
func (l *LoggerService) Info(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		l.Logger.Printf("[INFO] "+format, v...)
	}
}

// エラーレベルのログを出力します
func (l *LoggerService) Error(format string, v ...interface{}) {
	if l.level <= ErrorLevel {
		l.Logger.Printf("[ERROR] "+format, v...)
	}
}

// ロガーのログレベルを設定します
func (l *LoggerService) SetLevel(level LogLevel) {
	l.level = level
}

// ロガーの出力先を設定します
func (l *LoggerService) SetOutput(w *os.File) {
	l.Logger.SetOutput(w)
}

// ロガーのプレフィックスを設定します
func (l *LoggerService) SetPrefix(prefix string) {
	l.Logger.SetPrefix(prefix)
}

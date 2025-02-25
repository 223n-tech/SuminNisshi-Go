// internal/util/parse.go
// Package utilは、アプリケーション全体で使用されるユーティリティ関数を提供します

// Package util provides utility functions used throughout the application
package util

import (
	"fmt"
	"strconv"
	"time"
)

// 文字列をint64に変換する
func ParseInt64(s string) int64 {
    v, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0
    }
    return v
}

// "2006-01-02"形式の文字列をtime.Timeに変換する
func ParseDate(s string) time.Time {
    t, err := time.Parse("2006-01-02", s)
    if err != nil {
        return time.Time{}
    }
    return t
}

// "15:04" 形式の文字列をtime.Timeに変換する
func ParseTime(s string) time.Time {
    t, err := time.Parse("15:04", s)
    if err != nil {
        return time.Time{}
    }
    return t
}

// 時間の表示形式を整形
func formatDuration(duration float64) string {
	hours := int(duration)
	minutes := int((duration - float64(hours)) * 60)
	return fmt.Sprintf("%d時間%d分", hours, minutes)
}

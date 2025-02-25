// internal/models/sleep_diary.go
// sleep_diaryは、睡眠日誌の基本情報を管理する構造体を提供します。

// Package models provides data models for the application.
package models

import (
	"database/sql"
	"time"
)

/*
	睡眠日誌の基本情報を管理する構造体
*/
type SleepDiary struct {
	ID         int64          `db:"id"`
	UserID     int64          `db:"user_id"`
	StartDate  time.Time      `db:"start_date"`
	EndDate    time.Time      `db:"end_date"`
	DiaryName  string         `db:"diary_name"`
	Note       sql.NullString `db:"note"`
	Created    time.Time      `db:"created"`
	Modified   time.Time      `db:"modified"`
	Deleted    sql.NullTime   `db:"deleted"`
}

/*
	睡眠日誌のバリデーション
*/
func (d *SleepDiary) Validate() error {
	// TODO: 実装
	return nil
}

/*
	睡眠日誌の期間を計算
*/
func (d *SleepDiary) CalculateDuration() int {
	return int(d.EndDate.Sub(d.StartDate).Hours() / 24)
}

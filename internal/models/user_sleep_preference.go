// internal/models/user_sleep_preference.go
// user_sleep_preferenceは、ユーザーの睡眠設定を管理する構造体を提供します。

// Package models provides data models for the application.
package models

import (
	"database/sql"
	"time"
)

/*
	ユーザーの睡眠設定を管理する構造体
*/
type UserSleepPreference struct {
	ID                  int64        `db:"id"`
	UserID              int64        `db:"user_id"`
	PreferredBedtime    time.Time    `db:"preferred_bedtime"`
	PreferredWakeupTime time.Time    `db:"preferred_wakeup_time"`
	SleepGoalHours     int          `db:"sleep_goal_hours"`
	IsReminderEnabled  bool         `db:"is_reminder_enabled"`
	Created            time.Time    `db:"created"`
	Modified           time.Time    `db:"modified"`
	Deleted            sql.NullTime `db:"deleted"`
}

/*
	睡眠設定のバリデーション
*/
func (p *UserSleepPreference) Validate() error {
	// TODO: 実装
	return nil
}

/*
	目標睡眠時間を計算
*/
func (p *UserSleepPreference) CalculateTargetSleepDuration() time.Duration {
	return time.Duration(p.SleepGoalHours) * time.Hour
}

/*
	指定時刻が目標時間内かチェック
*/
func (p *UserSleepPreference) IsWithinTargetTime(t time.Time) bool {
	// 時刻のみを比較するために日付部分を無視
	targetTime := t.Format("15:04")
	bedTime := p.PreferredBedtime.Format("15:04")
	wakeTime := p.PreferredWakeupTime.Format("15:04")
	
	return targetTime >= bedTime || targetTime <= wakeTime
}

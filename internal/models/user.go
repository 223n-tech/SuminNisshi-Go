// internal/models/user.go
// userは、ユーザー情報を管理する構造体を提供します。

// Package models provides data models for the application.
package models

import (
	"database/sql"
	"time"
)

/*
	ユーザー情報を管理する構造体
*/
type User struct {
	ID                int64        `db:"id"`
	Email             string       `db:"email"`
	DisplayName       string       `db:"display_name"`
	PasswordHash      string       `db:"password_hash"`
	TimeZone          string       `db:"time_zone"`
	LastLoginDatetime sql.NullTime `db:"last_login_datetime"`
	Created           time.Time    `db:"created"`
	Modified          time.Time    `db:"modified"`
	Deleted           sql.NullTime `db:"deleted"`
}

/*
	通知設定の構造体
	既存の構造体は残しつつ、新しい要件に合わせて拡張
*/
type NotificationSettings struct {
	EmailEnabled    bool
	BedtimeReminder bool
	ReminderTime    string
	WeeklyReport    bool
}

/*
	ユーザー情報のバリデーション
*/
func (u *User) Validate() error {
	// TODO: 実装
	return nil
}

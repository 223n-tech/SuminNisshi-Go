// internal/models/user.go
// ユーザー情報の構造体を定義する
package models

/*
	User ユーザー情報の構造体
*/
type User struct {
	Name     string
	Email    string
	Timezone string
}

/*
	NotificationSettings 通知設定の構造体
*/
type NotificationSettings struct {
	EmailEnabled    bool
	BedtimeReminder bool
	ReminderTime    string
	WeeklyReport    bool
}

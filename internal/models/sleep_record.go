// internal/models/sleep_record.go
// sleep_recordは、睡眠記録の詳細を管理する構造体を提供します。

// Package models provides data models for the application.
package models

import (
	"database/sql"
	"time"
)

/*
	睡眠記録の詳細を管理する構造体
*/
type SleepRecord struct {
	ID           int64          `db:"id"`
	SleepDiaryID int64          `db:"sleep_diary_id"`
	SleepStateID int64          `db:"sleep_state_id"`
	RecordDate   time.Time      `db:"record_date"`
	TimeSlot     time.Time      `db:"time_slot"`
	RecordType   string         `db:"record_type"` // ENUM: STATE, EVENT, MEAL
	MealTypeID   sql.NullInt64  `db:"meal_type_id"`
	Note         sql.NullString `db:"note"`
	Created      time.Time      `db:"created"`
	Modified     time.Time      `db:"modified"`
	Deleted      sql.NullTime   `db:"deleted"`
}

/*
	記録種別を定義する定数
*/
const (
	RecordTypeState = "STATE"
	RecordTypeEvent = "EVENT"
	RecordTypeMeal  = "MEAL"
)

/*
	関連データを含む睡眠記録
*/
type SleepRecordWithRelations struct {
	SleepRecord
	State    SleepState `db:"state"`
	MealType *MealType  `db:"meal_type"`
	Diary    SleepDiary `db:"diary"`
}

/*
	睡眠記録のバリデーション
*/
func (r *SleepRecord) Validate() error {
	// TODO: 実装
	return nil
}

/*
	時間枠が30分単位かチェック
*/
func (r *SleepRecord) IsValidTimeSlot() bool {
	minutes := r.TimeSlot.Minute()
	return minutes == 0 || minutes == 30
}

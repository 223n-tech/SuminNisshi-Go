// internal/models/sleep_state.go
// sleep_stateは、睡眠状態のマスターデータを管理する構造体を提供します。

// Package models provides data models for the application.
package models

import (
	"database/sql"
	"time"
)

/*
	睡眠状態のマスターデータを管理する構造体
*/
type SleepState struct {
	ID               int64          `db:"id"`
	StateName        string         `db:"state_name"`
	StateCode        string         `db:"state_code"`
	StateDescription sql.NullString `db:"state_description"`
	DisplaySymbol    string         `db:"display_symbol"`
	DisplayOrder     int            `db:"display_order"`
	Created          time.Time      `db:"created"`
	Modified         time.Time      `db:"modified"`
	Deleted          sql.NullTime   `db:"deleted"`
}

/*
	睡眠状態コードを定義する定数
*/
const (
	StateCodeSleeping   = "SLEEPING"
	StateCodeAwakeInBed = "AWAKE_IN_BED"
	StateCodeAwake      = "AWAKE"
	StateCodeDrowsiness = "DROWSINESS"
	StateCodeMedication = "MEDICATION"
)

/*
	デフォルトの睡眠状態を返す
*/
func DefaultSleepStates() []SleepState {
	return []SleepState{
		{StateName: "睡眠中", StateCode: StateCodeSleeping, DisplaySymbol: "■", DisplayOrder: 1},
		{StateName: "床で覚醒", StateCode: StateCodeAwakeInBed, DisplaySymbol: "╱", DisplayOrder: 2},
		{StateName: "通常覚醒", StateCode: StateCodeAwake, DisplaySymbol: "□", DisplayOrder: 3},
		{StateName: "強い眠気", StateCode: StateCodeDrowsiness, DisplaySymbol: "Z", DisplayOrder: 4},
		{StateName: "睡眠薬服用", StateCode: StateCodeMedication, DisplaySymbol: "×", DisplayOrder: 5},
	}
}

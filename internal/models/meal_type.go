// internal/models/meal_type.go
// meal_typeは、食事種別のマスターデータを管理する構造体を提供します。

// Package models provides data models for the application.
package models

import (
	"database/sql"
	"time"
)

/*
	食事種別のマスターデータを管理する構造体
*/
type MealType struct {
	ID            int64        `db:"id"`
	TypeName      string       `db:"type_name"`
	TypeCode      string       `db:"type_code"`
	DisplaySymbol string       `db:"display_symbol"`
	DisplayOrder  int          `db:"display_order"`
	Created       time.Time    `db:"created"`
	Modified      time.Time    `db:"modified"`
	Deleted       sql.NullTime `db:"deleted"`
}

/*
	食事種別コードを定義する定数
*/
const (
	MealCodeBreakfast = "BREAKFAST"
	MealCodeLunch     = "LUNCH"
	MealCodeDinner    = "DINNER"
	MealCodeSnack     = "SNACK"
)

/*
	デフォルトの食事種別を返す
*/
func DefaultMealTypes() []MealType {
	return []MealType{
		{TypeName: "朝食", TypeCode: MealCodeBreakfast, DisplaySymbol: "▲", DisplayOrder: 1},
		{TypeName: "昼食", TypeCode: MealCodeLunch, DisplaySymbol: "●", DisplayOrder: 2},
		{TypeName: "夕食", TypeCode: MealCodeDinner, DisplaySymbol: "■", DisplayOrder: 3},
		{TypeName: "軽食", TypeCode: MealCodeSnack, DisplaySymbol: "○", DisplayOrder: 4},
	}
}

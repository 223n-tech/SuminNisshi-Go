# パッケージ: models

## 1. 概要

Package models provides data models for the application.
internal/models/sleep_diary.go
sleep_diaryは、睡眠日誌の基本情報を管理する構造体を提供します。

## 2. 型

### 2-1. `MealType`

食事種別のマスターデータを管理する構造体

#### 2-1-1. フィールド

- `ID`: `int64`

- `TypeName`: `string`

- `TypeCode`: `string`

- `DisplaySymbol`: `string`

- `DisplayOrder`: `int`

- `Created`: `time.Time`

- `Modified`: `time.Time`

- `Deleted`: `sql.NullTime`

### 2-2. `NotificationSettings`

通知設定の構造体
既存の構造体は残しつつ、新しい要件に合わせて拡張

#### 2-2-1. フィールド

- `EmailEnabled`: `bool`

- `BedtimeReminder`: `bool`

- `ReminderTime`: `string`

- `WeeklyReport`: `bool`

### 2-3. `PDFExport`

PDF出力データ

#### 2-3-1. フィールド

- `User`: `*User`

- `DiaryName`: `string`

- `Records`: `[]*SleepRecord`

- `Statistics`: `*SleepStatistics`

- `GeneratedAt`: `time.Time`

- `Options`: `*PDFExportOptions`

### 2-4. `PDFExportData`

PDF出力用のデータを管理する構造体

#### 2-4-1. フィールド

- `User`: `User`

- `SleepDiary`: `SleepDiary`

- `Records`: `[]*SleepRecord`

- `States`: `*ast.MapType`

- `MealTypes`: `*ast.MapType`

- `Preferences`: `UserSleepPreference`

### 2-5. `PDFExportOptions`

PDF出力設定を管理する構造体

#### 2-5-1. フィールド

- `StartDate`: `time.Time`

- `EndDate`: `time.Time`

- `IncludeNote`: `bool`

- `Quality`: `string`

- `Language`: `string`

### 2-6. `PDFSleepRecord`

PDF出力用の睡眠記録データ

#### 2-6-1. フィールド

- `Date`: `time.Time`

- `BedTime`: `string`

- `WakeTime`: `string`

- `Duration`: `float64`

- `Score`: `int`

### 2-7. `PDFStatistics`

PDF出力用の統計データ

#### 2-7-1. フィールド

- `AverageDuration`: `float64`

- `AverageBedTime`: `string`

- `AverageWakeTime`: `string`

- `AverageScore`: `float64`

- `StartDate`: `time.Time`

- `EndDate`: `time.Time`

- `TotalDays`: `int`

### 2-8. `PDFTemplate`

PDFテンプレート用の設定

#### 2-8-1. フィールド

- `FontPath`: `string`

- `PageWidth`: `float64`

- `PageHeight`: `float64`

- `Margin`: `float64`

### 2-9. `PDFTimeSlot`

時間枠の定義

#### 2-9-1. フィールド

- `Hour`: `int`

- `Symbol`: `string`

- `IsAwake`: `bool`

- `HasEvent`: `bool`

### 2-10. `SleepDiary`

睡眠日誌の基本情報を管理する構造体

#### 2-10-1. フィールド

- `ID`: `int64`

- `UserID`: `int64`

- `StartDate`: `time.Time`

- `EndDate`: `time.Time`

- `DiaryName`: `string`

- `Note`: `sql.NullString`

- `Created`: `time.Time`

- `Modified`: `time.Time`

- `Deleted`: `sql.NullTime`

### 2-11. `SleepRecord`

睡眠記録の詳細を管理する構造体

#### 2-11-1. フィールド

- `ID`: `int64`

- `SleepDiaryID`: `int64`

- `SleepStateID`: `int64`

- `RecordDate`: `time.Time`

- `TimeSlot`: `time.Time`

- `RecordType`: `string`

- `MealTypeID`: `sql.NullInt64`

- `Note`: `sql.NullString`

- `Created`: `time.Time`

- `Modified`: `time.Time`

- `Deleted`: `sql.NullTime`

### 2-12. `SleepRecordWithRelations`

関連データを含む睡眠記録

#### 2-12-1. フィールド

- `State`: `SleepState`

- `MealType`: `*MealType`

- `Diary`: `SleepDiary`

### 2-13. `SleepState`

睡眠状態のマスターデータを管理する構造体

#### 2-13-1. フィールド

- `ID`: `int64`

- `StateName`: `string`

- `StateCode`: `string`

- `StateDescription`: `sql.NullString`

- `DisplaySymbol`: `string`

- `DisplayOrder`: `int`

- `Created`: `time.Time`

- `Modified`: `time.Time`

- `Deleted`: `sql.NullTime`

### 2-14. `SleepStatistics`

睡眠統計データ

#### 2-14-1. フィールド

- `AverageSleepHours`: `float64`

- `AverageSleepQuality`: `float64`

- `TargetAchievementRate`: `float64`

- `MostCommonBedtime`: `time.Time`

- `MostCommonWaketime`: `time.Time`

### 2-15. `User`

ユーザー情報を管理する構造体

#### 2-15-1. フィールド

- `ID`: `int64`

- `Email`: `string`

- `DisplayName`: `string`

- `PasswordHash`: `string`

- `TimeZone`: `string`

- `LastLoginDatetime`: `sql.NullTime`

- `Created`: `time.Time`

- `Modified`: `time.Time`

- `Deleted`: `sql.NullTime`

### 2-16. `UserSleepPreference`

ユーザーの睡眠設定を管理する構造体

#### 2-16-1. フィールド

- `ID`: `int64`

- `UserID`: `int64`

- `PreferredBedtime`: `time.Time`

- `PreferredWakeupTime`: `time.Time`

- `SleepGoalHours`: `int`

- `IsReminderEnabled`: `bool`

- `Created`: `time.Time`

- `Modified`: `time.Time`

- `Deleted`: `sql.NullTime`

## 3. 関数

なし

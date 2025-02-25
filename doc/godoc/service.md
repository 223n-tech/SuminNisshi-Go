# パッケージ: service

## 1. 概要

Package service provides application services.
internal/service/user_service.go
user_serviceは、ユーザー関連のサービスを提供します。

## 2. 型

### 2-1. `DashboardStats`

ダッシュボードに表示するデータ

#### 2-1-1. フィールド

- `TotalSleepHours`: `float64`

- `AverageSleepHours`: `float64`

- `SleepQualityScore`: `int`

- `TargetAchievement`: `int`

- `RecentRecords`: `[]models.SleepRecord`

### 2-2. `EmailService`

メール送信サービス

#### 2-2-1. フィールド

- `s`: `*Service`

### 2-3. `LogLevel`

ログのレベルを表す型

### 2-4. `LoggerService`

標準ログパッケージを拡張したロガー

#### 2-4-1. フィールド

- `level`: `LogLevel`

### 2-5. `PDFService`

PDF出力関連のサービス

#### 2-5-1. フィールド

- `s`: `*Service`

### 2-6. `Service`

アプリケーションのサービス層を表す構造体

#### 2-6-1. フィールド

- `repo`: `repository.Repository`

- `logger`: `*LoggerService`

- `user`: `*UserService`

- `diary`: `*SleepDiaryService`

- `record`: `*SleepRecordService`

- `pdf`: `*PDFService`

- `email`: `*EmailService`

### 2-7. `SleepDiaryService`

睡眠日誌関連のサービス

#### 2-7-1. フィールド

- `s`: `*Service`

### 2-8. `SleepRecordFilter`

フィルター条件

#### 2-8-1. フィールド

- `StartDate`: `string`

- `EndDate`: `string`

- `StateID`: `int64`

### 2-9. `SleepRecordService`

睡眠記録関連のサービス

#### 2-9-1. フィールド

- `s`: `*Service`

### 2-10. `UserService`

ユーザー関連のサービス

#### 2-10-1. フィールド

- `s`: `*Service`

## 3. 関数

なし

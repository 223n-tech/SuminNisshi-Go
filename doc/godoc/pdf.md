# パッケージ: pdf

## 1. 概要

Package pdf provides functionality to generate PDFs.
internal/pdf/pdf.go
pdfは、PDFを生成するための機能を提供します。

## 2. 型

### 2-1. `Generator`

PDFを生成するための構造体

#### 2-1-1. フィールド

- `pdf`: `*gopdf.GoPdf`

- `fontPath`: `string`

### 2-2. `SleepRecord`

個別の睡眠記録

#### 2-2-1. フィールド

- `Date`: `time.Time`

- `BedTime`: `string`

- `WakeTime`: `string`

- `Duration`: `float64`

- `Score`: `int`

### 2-3. `SleepRecordData`

PDFに出力する睡眠記録データ

#### 2-3-1. フィールド

- `StartDate`: `time.Time`

- `EndDate`: `time.Time`

- `TotalDays`: `int`

- `Records`: `[]SleepRecord`

- `AverageDuration`: `float64`

- `AverageBedTime`: `string`

- `AverageWakeTime`: `string`

- `AverageScore`: `float64`

## 3. 関数

なし

# パッケージ: util

## 1. 概要

Package util provides utility functions used throughout the application
internal/util/parse.go
Package utilは、アプリケーション全体で使用されるユーティリティ関数を提供します

## 2. 型

なし

## 3. 関数

### 3-1. `ParseDate`

"2006-01-02"形式の文字列をtime.Timeに変換する

### 3-2. `ParseInt64`

文字列をint64に変換する

### 3-3. `ParseTime`

"15:04" 形式の文字列をtime.Timeに変換する

### 3-4. `formatDuration`

時間の表示形式を整形

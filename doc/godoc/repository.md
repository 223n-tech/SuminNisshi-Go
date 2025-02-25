# パッケージ: repository

## 1. 概要

Package repository provides interfaces for abstracting database operations.
This package provides interfaces for performing database operations.
This allows you to minimize code changes when changing database implementations.
internal/repository/repository.go
repositoryは、データベースの操作を抽象化するためのインターフェイスを提供します。
このパッケージは、データベースの操作を行うためのインターフェイスを提供します。
これにより、データベースの実装を変更する際に、コードの変更を最小限に抑えることができます。

## 2. 型

### 2-1. `MealTypeRepository`

食事種別のリポジトリーインターフェイス

### 2-2. `Repository`

全リポジトリのインターフェイス

### 2-3. `SleepDiaryRepository`

睡眠日誌のリポジトリーインターフェイス

### 2-4. `SleepRecordRepository`

睡眠記録のリポジトリーインターフェイス

### 2-5. `SleepStateRepository`

睡眠状態のリポジトリーインターフェイス

### 2-6. `UserRepository`

ユーザー情報のリポジトリーインターフェイス

### 2-7. `UserSleepPreferenceRepository`

ユーザー睡眠設定のリポジトリーインターフェイス

## 3. 関数

なし

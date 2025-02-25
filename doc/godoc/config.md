# パッケージ: config

## 1. 概要

Package config provides a structure to hold the application-wide configuration.

## 2. 型

### 2-1. `Config`

アプリケーション全体の設定を保持する構造体

#### 2-1-1. フィールド

- `Server`: `ServerConfig`

- `Database`: `DatabaseConfig`

### 2-2. `DatabaseConfig`

データベース関連の設定

#### 2-2-1. フィールド

- `Host`: `string`

- `Port`: `int`

- `User`: `string`

- `Password`: `string`

- `DBName`: `string`

### 2-3. `ServerConfig`

サーバー関連の設定

#### 2-3-1. フィールド

- `Port`: `int`

- `Host`: `string`

- `BaseURL`: `string`

## 3. 関数

### 3-1. `getEnvInt`

環境変数から整数を取得

### 3-2. `getEnvStr`

環境変数から文字列を取得

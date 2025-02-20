# SuiminNisshi-Go プロジェクト概要

## 基本情報

* リポジトリ: [https://github.com/223n-tech/SuiminNisshi-Go](https://github.com/223n-tech/SuiminNisshi-Go)
* アプリケーション名: SuiminNisshi-Go
* Goバージョン: 1.22.1
* OS: debian bookworm-slim

## 開発環境

* devcontainer + VS Code使用
* MariaDB最新版
* AdminLTE最新版を使用したUI

## 現在の開発状況

### 作成済みページ/機能

1. 認証関連
   * ログインページ
   * ユーザー登録ページ
   * パスワードリセットページ

2. メイン機能
   * ダッシュボード
   * 睡眠記録一覧・詳細・フォーム
   * 統計情報
   * プロフィール管理
   * 設定画面
   * 利用規約ページ
   * プライバシーポリシーページ

3. その他
   * エラーページ（403, 404, 500）

### 作成中のページ/機能

1. アカウント削除確認ページ
2. データエクスポート確認ページ

### 未実装の機能

1. データベース実装
2. メール送信機能
3. バリデーション
4. テスト
5. ログ管理
6. セッション管理

## 技術スタック

1. フレームワーク/ライブラリ
   * github.com/go-chi/chi/v5（ルーティング）
   * github.com/go-chi/cors（CORS処理）
   * github.com/go-sql-driver/mysql（データベース）

2. アーキテクチャ
   * リポジトリパターン採用
   * ミドルウェアによる認証処理
   * テンプレートエンジンによるレンダリング

## ディレクトリ構造の主要部分

```bash
.
├── cmd/
│   └── suiminnisshi/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handler/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   └── service/
├── web/
│   ├── views/
│   │   ├── errors/
│   │   ├── layouts/
│   │   ├── pages/
│   │   └── partials/
│   └── static/
│       └── adminlte/
└── go.mod
```

## 現在の課題

1. modelsパッケージの実装
   * User構造体の実装
   * NotificationSettings構造体の実装
   * パッケージのインポートパス修正

2. テンプレート関連の修正
   * TemplateDataのMeta対応
   * エラーページのテンプレート修正
   * 各ページのデータ構造の整理

3. フォーム処理の改善
   * バリデーション実装
   * エラーハンドリング
   * フラッシュメッセージ対応

## セットアップ手順

1. devcontainer環境の構築
2. AdminLTEのセットアップ
3. 環境変数の設定（direnv使用）
4. データベースのセットアップ（未実装）

## 次のステップ

1. modelsパッケージの完全実装
2. データベース接続とリポジトリの実装
3. バリデーション機能の実装
4. テストの追加
5. ログ管理の実装
6. セッション管理の実装

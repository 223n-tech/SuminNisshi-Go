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
   * アカウント削除確認ページ
   * データエクスポート確認ページ

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
│   │  └── config.go
│   ├── handler/
│   │  ├── account_deletion.go
│   │  ├── auth.go
│   │  ├── dashboard.go
│   │  ├── error.go
│   │  ├── password_reset.go
│   │  ├── privacy.go
│   │  ├── profile.go
│   │  ├── register.go
│   │  ├── router.go
│   │  ├── settings.go
│   │  ├── sleep_records.go
│   │  ├── statistics.go
│   │  ├── template.go
│   │  └── terms.go
│   ├── middleware/
│   │  └── security.go
│   ├── models/
│   │  ├── meal_type.go
│   │  ├── pdf_export.go
│   │  ├── sleep_diary.go
│   │  ├── sleep_record.go
│   │  ├── sleep_state.go
│   │  ├── user_sleep_preference.go
│   │  └── user.go
│   ├── pdf/
│   │  └── pdf.go
│   ├── repository/
│   │  ├── mysql/
│   │  │  ├── db.go
│   │  │  ├── meal_type_repository.go
│   │  │  ├── mysql_repository.go
│   │  │  ├── sleep_diary_repository.go
│   │  │  ├── sleep_record_repository.go
│   │  │  ├── sleep_state_repository.go
│   │  │  ├── user_repository.go
│   │  │  └── user_sleep_preference_repository.go
│   │  └── repository.go
│   ├── service/
│   │  ├── logger_service.go
│   │  ├── mail_service.go
│   │  ├── pdf_service.go
│   │  ├── service.go
│   │  ├── sleep_diary_service.go
│   │  ├── sleep_record_service.go
│   │  └── user_service.go
│   └── util/
│       └── parse.go
├── static/
├── tools/
│   └── doc-template-generator.go
├── web/
│   ├── static/
│   │   └── adminlte/
│   ├── template/
│   │   ├── charts/
│   │   ├── examples/
│   │   ├── forms/
│   │   ├── layout/
│   │   ├── mailbox/
│   │   ├── search/
│   │   ├── tables/
│   │   └── UI/
│   └── views/
│        ├── errors/
│        │  ├── 403.html
│        │  ├── 404.html
│        │  └── 500.html
│        ├── layouts/
│        │  └── base.html
│        ├── pages/
│        │  ├── account-deletion.html
│        │  ├── dashboard.html
│        │  ├── delete-account.html
│        │  ├── export-data.html
│        │  ├── forgot-password.html
│        │  ├── login.html
│        │  ├── privacy.html
│        │  ├── profile.html
│        │  ├── register.html
│        │  ├── reset-password.html
│        │  ├── settings.html
│        │  ├── sleep-records-detail.html
│        │  ├── sleep-records-form.html
│        │  ├── sleep-records.html
│        │  ├── statistics.html
│        │  └── terms.html
│        └── partials/
│            ├── footer.html
│            ├── navbar.html
│            └── sidebar.html
└── go.mod
```

## セットアップ手順

1. devcontainer環境の構築
2. AdminLTEのセットアップ
3. 環境変数の設定（direnv使用）
4. データベースのセットアップ（未実装）

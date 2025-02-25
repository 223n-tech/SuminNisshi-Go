# パッケージ: handler

## 1. 概要

Package handler provides HTTP handlers for the application.

## 2. 型

### 2-1. `AccountDeletionHandler`

アカウント削除画面のハンドラー

#### 2-1-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

### 2-2. `AuthHandler`

認証関連のハンドラー

#### 2-2-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

### 2-3. `DashboardHandler`

ダッシュボード関連のハンドラー

#### 2-3-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

### 2-4. `ErrorHandler`

エラーページのハンドラー

#### 2-4-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

- `logger`: `*log.Logger`

### 2-5. `Flash`

フラッシュメッセージの構造体

#### 2-5-1. フィールド

- `Type`: `string`

- `Message`: `string`

### 2-6. `PasswordResetHandler`

パスワードリセット関連のハンドラー

#### 2-6-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

### 2-7. `PrivacyHandler`

プライバシーポリシーページのハンドラー

#### 2-7-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

### 2-8. `ProfileHandler`

プロフィール関連のハンドラー

#### 2-8-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

### 2-9. `ProfileUpdateRequest`

プロフィール更新リクエストの構造体

#### 2-9-1. フィールド

- `DisplayName`: `string`

- `Email`: `string`

- `CurrentPassword`: `string`

- `NewPassword`: `string`

- `PasswordConfirm`: `string`

- `Timezone`: `string`

### 2-10. `RegisterData`

新規登録画面のデータ

#### 2-10-1. フィールド

- `Name`: `string`

- `Email`: `string`

- `Password`: `string`

- `PasswordConfirmation`: `string`

- `Terms`: `bool`

- `Error`: `string`

### 2-11. `RegisterHandler`

新規登録画面のハンドラー

#### 2-11-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

### 2-12. `RouterWrapper`

chi.Routerのラッパーです

### 2-13. `SettingsHandler`

設定ページのハンドラー

#### 2-13-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

### 2-14. `SleepRecordHandler`

睡眠記録関連のハンドラー

#### 2-14-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

### 2-15. `StatisticsHandler`

統計情報関連のハンドラー

#### 2-15-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

### 2-16. `TemplateData`

テンプレートに渡すデータの構造体

#### 2-16-1. フィールド

- `Title`: `string`

- `ActiveMenu`: `string`

- `User`: `*models.User`

- `Data`: `*ast.MapType`

- `Flash`: `*Flash`

- `Meta`: `*ast.MapType`

### 2-17. `TemplateManager`

テンプレートを管理する構造体

#### 2-17-1. フィールド

- `templates`: `*ast.MapType`

- `mutex`: `sync.RWMutex`

- `basePath`: `string`

- `funcMap`: `template.FuncMap`

- `embedFS`: `*embed.FS`

- `standalonePages`: `[]string`

- `logger`: `*log.Logger`

- `service`: `*service.Service`

### 2-18. `TermsHandler`

利用規約画面のハンドラー

#### 2-18-1. フィールド

- `templates`: `*TemplateManager`

- `service`: `*service.Service`

### 2-19. `userContextKey`

コンテキストのキー

## 3. 関数

### 3-1. `GetUserFromContext`

コンテキストからユーザー情報を取得

### 3-2. `GetUserIDFromContext`

コンテキストからユーザーIDを取得

### 3-3. `RequireAuth`

認証が必要なリクエストに対してミドルウェアを適用

### 3-4. `getScoreColorClass`

睡眠スコアに応じた色のクラスを取得

### 3-5. `makeTemplateFuncMap`

テンプレート関数を作成

### 3-6. `parseInt`

文字列を整数に変換

### 3-7. `validateSleepRecord`

睡眠記録のバリデーション

### 3-8. `validateTimeRange`

時間範囲の妥当性をチェック

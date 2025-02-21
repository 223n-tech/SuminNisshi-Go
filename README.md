# SuminNisshi-Go

簡単に言えば、[SuiminNisshi](https://github.com/223n-tech/SuminNisshi)のGo言語版です。

## はじめに

このREADME内の各種コマンドは、devcontainer環境下で実行した場合です。

## 注意

まだ、ページは **ハリボテの見た目** しか作成していません。
具体的な実装は、未実装なので注意してください。

## ページ

| handler             | アドレス                                                                  | ページ名                     | 役割 | モック実装 | 処理実装 |
| ------------------- | ------------------------------------------------------------------------- | ---------------------------- | ---- | :--------: | :------: |
| N/A                 | [/](http://localhost:8080/)                                               | トップページ                 |      |     x      |    x     |
| account_deletion.go | [/account/delete](http://localhost:8080/account/delete)                   | アカウント削除確認ページ     |      |     o      |    x     |
| dashboard.go        | [/dashboard](http://localhost:8080/dashboard)                             | ダッシュボード               |      |     o      |    x     |
| auth.go             | [/login](http://localhost:8080/login)                                     | ログインページ               |      |     o      |    x     |
| auth.go             | [/logout](http://localhost:8080/logout)                                   | ログアウトページ             |      |     o      |    x     |
| error.go            | [/{存在しないページ}](http://localhost:8080/abc)                          | 404ページ                    |      |     o      |    x     |
| error.go            | [未設定](http://localhost:8080/)                                          | 403ページ                    |      |     o      |    x     |
| error.go            | [未設定](http://localhost:8080/)                                          | 405ページ                    |      |     o      |    x     |
| error.go            | [未設定](http://localhost:8080/)                                          | 500ページ                    |      |     o      |    x     |
| password_reset.go   | [/forgot-password](http://localhost:8080/forgot-password)                 | パスワード忘れページ         |      |     o      |    x     |
| password_reset.go   | [/reset-password/{token}](http://localhost:8080/reset-password/abc)       | パスワード再設定ページ       |      |     o      |    x     |
| privacy.go          | [/privacy](http://localhost:8080/privacy)                                 | プライバシーポリシーページ   |      |     o      |    x     |
| profile.go          | [/profile](http://localhost:8080/profile)                                 | プロフィールページ           |      |     o      |    x     |
| register.go         | [/register](http://localhost:8080/register)                               | 新規アカウント登録ページ     |      |     x      |    x     |
| settings.go         | [/settings](http://localhost:8080/settings)                               | 設定ページ                   |      |     o      |    x     |
| settings.go         | [/settings/export/csv](http://localhost:8080/settings/export/csv)         | 設定ページ（CSV出力）        |      |     o      |    x     |
| settings.go         | [/settings/export/json](http://localhost:8080/settings/export/json)       | 設定ページ（JSON出力）       |      |     o      |    x     |
| settings.go         | [/settings/account/delete](http://localhost:8080/settings/account/delete) | 設定ページ（アカウント削除） |      |     o      |    x     |
| sleep_records.go    | [/sleep-records/](http://localhost:8080/sleep-records)                    | 睡眠記録一覧ページ           |      |     x      |    x     |
| sleep_records.go    | [/sleep-records/new](http://localhost:8080/sleep-records/new)             | 睡眠記録入力ページ           |      |     x      |    x     |
| sleep_records.go    | [/sleep-records/{id}](http://localhost:8080/sleep-records/1)              | 睡眠記録詳細ページ           |      |     x      |    x     |
| sleep_records.go    | [/sleep-records/{id}/edit](http://localhost:8080/sleep-records/1/edit)    | 睡眠記録編集ページ           |      |     x      |    x     |
| sleep_records.go    | [/api/sleep-records/](http://localhost:8080/sleep-records/1/edit)         | (API)睡眠記録一覧            |      |     x      |    x     |
| statistics.go       | [/statistics](http://localhost:8080/statistics)                           | 統計情報ページ               |      |     x      |    x     |
| statistics.go       | [/statistics/data](http://localhost:8080/statistics/data)                 | 統計情報ページ               |      |     x      |    x     |
| terms.go            | [/terms](http://localhost:8080/terms)                                     | 利用規約ページ               |      |     x      |    x     |

## 基本コマンド

### ホストの起動

```bash
go run cmd/suiminnisshi/main.go
```

### ホストの終了

ホストの起動中に、Ctrl+Cで終了させることができます。

### AdminLTEの再セットアップ

```bash
bash /usr/local/bin/setup-adminlte.sh
```

## devcontainer環境

### APP

* image: [debian:bookworm-slim](https://hub.docker.com/_/debian)
* OS: [Debian bookworm](https://www.debian.org/releases/bookworm/)
* 開発言語: [Go (Ver.1.22.1)](https://tip.golang.org/doc/devel/release)
* 環境変数管理: [direnv](https://github.com/direnv/direnv)
* テンプレート: [AdminLTE](https://adminlte.io/)

### DB

* image: [mariadb:latest](https://hub.docker.com/_/mariadb)
* DB: [MariaDB](https://mariadb.org/)
* 設定
  * DB_HOST = db
  * USER_DB = suiminnisshi
  * DB_PASSWORD = suiminnisshi_password
  * DB_NAME = suiminnisshi

### ポート転送

* 8080: アプリケーションポート
* 3306: MariaDBポート

## 開発メモ

### godocについて

godocコマンドで、ドキュメントページにアクセスすることが可能です。
internalディレクトリ以下のモジュールは、対象外です。

```bash
godoc
```

### .envrcファイルについて

* `.envrc`ファイルが存在しない場合、devcontainerビルド時にテンプレートファイルから自動生成されます。
  * テンプレートファイルは、 `.devcontainer/.envrc.template` です。
* すでに`.envrc`ファイルが存在する場合、自動生成されません。

### ターミナルを起動するとdirenvのエラーが表示される

`.envrc`ファイルの読み込みを有効にしてください。

```bash
direnv allow
```

### 新しくページを追加したい

1. `web/template/pages`ディレクトリに作成したいページ名のhtmlファイルを作成する。
2. ページは基本レイアウトを継承して必要なコンテンツブロックを実装します。

    ```html
    {{define "content"}}
      <!-- ページコンテンツをここに記述 -->
    {{end}}
    ```

3. AdminLTEを使う場合には、
    * 必要なCSSファイルをstyleブロックで追加します。
    * 必要なJavaScriptファイルをscriptブロックで追加します。

4. main.goにハンドラーの初期化と登録（必要な場合にのみ実施）
    * 必要な場合には、main.goにハンドラーを記載します。
    * 自動読み込みで良い場合には、不要です。

    ```go
      // ex. 設定例
      // 設定ハンドラーの初期化と登録
      settingsHandler := handler.NewSettingsHandler(tm)
      settingsHandler.RegisterRoutes(router)
    ```

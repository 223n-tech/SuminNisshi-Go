# SuminNisshi-Go

簡単に言えば、[SuiminNisshi](https://github.com/223n-tech/SuminNisshi)のGo言語版です。

## 1. はじめに

このREADME内の各種コマンドは、devcontainer環境下で実行した場合です。

## 2. 注意

まだ、ページは **ハリボテの見た目** しか作成していません。
具体的な実装は、未実装なので注意してください。

## 3. ページ

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

### 3-1. ルート設定について

パスルートの設定は、[internal/handler](./internal/handler/)にある各ハンドラー内で定義している`RegisterRoutes`関数で設定しています。

## 4. 基本コマンド

`make`コマンドで実行している詳細については、[Makefileファイル](./Makefile)を参照してください。

### 4-1. ホストの起動

```bash
make run
```

### 4-2. ホストの終了

ホストの起動中に、Ctrl+Cで終了させることができます。

### 4-3. AdminLTEの再セットアップ

```bash
make mod.install.adminlte
```

### 4-4. ツールのインストール

```bash
make tools.get
```

### 4-5. ドキュメントの生成

```bash
make tools.doc
```

### 4-6. lintの実行

```bash
make lint
```

## 5. devcontainer環境

### 5-1. APP

* image: [debian:bookworm-slim](https://hub.docker.com/_/debian)
* OS: [Debian bookworm](https://www.debian.org/releases/bookworm/)
* 開発言語: [Go (Ver.1.22.1)](https://tip.golang.org/doc/devel/release)
* 環境変数管理: [direnv](https://github.com/direnv/direnv)
* テンプレート: [AdminLTE](https://adminlte.io/)

### 5-2. DB

* image: [mariadb:latest](https://hub.docker.com/_/mariadb)
* DB: [MariaDB](https://mariadb.org/)
* 設定
  * DB_HOST = db
  * USER_DB = suiminnisshi
  * DB_PASSWORD = suiminnisshi_password
  * DB_NAME = suiminnisshi

### 5-3. ポート転送

* 8080: アプリケーションポート
* 3306: MariaDBポート

## 6. 開発メモ

### 6-1. ドキュメントについて

makeコマンドで、ドキュメントファイルの自動生成が可能です。
生成したファイルは、[doc/godocディレクトリ](./doc/godoc/)に保存されています。

```bash
make tools.doc
```

* 生成ファイル
  * [cmd.md](./doc/godoc/cmd.md)
  * [config.md](./doc/godoc/config.md)
  * [handler.md](./doc/godoc/handler.md)
  * [middleware.md](./doc/godoc/middleware.md)
  * [models.md](./doc/godoc/models.md)
  * [pdf.md](./doc/godoc/pdf.md)
  * [repository.md](./doc/godoc/repository.md)
  * [service.md](./doc/godoc/service.md)
  * [tools.md](./doc/godoc/tools.md)
  * [util.md](./doc/godoc/util.md)
* スクリプト
  * [doc-template-generator.go](./tools/doc-template-generator.go)
* テンプレートファイル
  * [package-template.md](./doc/godoc/package-template.md)

### 6-2. .envrcファイルについて

* ルートフォルダーに`.envrc`ファイルが存在しない場合、devcontainerビルド時にテンプレートファイルから自動生成されます。
  * テンプレートファイルは、 [.devcontainer/.envrc.template](./.devcontainer/.envrc.template) です。
* すでに、ルートフォルダーに`.envrc`ファイルが存在する場合、自動生成されません。

### 6-3. ターミナルを起動するとdirenvのエラーが表示される

`.envrc`ファイルの読み込みを有効にしてください。

```bash
direnv allow
```

### 6-4. 新しくページを追加したい

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

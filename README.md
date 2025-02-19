# SuminNisshi-Go

簡単に言えば、[SuiminNisshi](https://github.com/223n-tech/SuminNisshi)のGo言語版です。

## はじめに

このREADME内の各種コマンドは、devcontainer環境下で実行した場合です。

## 注意

まだ、ページは **ハリボテの見た目** しか作成していません。
具体的な実装は、未実装なので注意してください。

## ページ

| アドレス                                                               | ページ名                 | 役割 | モック実装 | 処理実装 |
| ---------------------------------------------------------------------- | ------------------------ | ---- | :--------: | :------: |
| [/](http://localhost:8080/)                                            | トップページ             |      |     x      |    x     |
| [/login](http://localhost:8080/login)                                  | ログインページ           |      |     o      |    x     |
| [/dashboard](http://localhost:8080/dashboard)                          | ダッシュボード           |      |     o      |    x     |
| [/settings](http://localhost:8080/settings)                            | 設定ページ               |      |     o      |    x     |
| [/profile](http://localhost:8080/profile)                              | プロフィールページ       |      |     o      |    x     |
| [/register](http://localhost:8080/register)                            | 新規アカウント登録ページ |      |     x      |    x     |
| [/forgot-password](http://localhost:8080/forgot-password)              | パスワード忘れページ     |      |     x      |    x     |
| [/sleep-records/](http://localhost:8080/sleep-records)                 | 睡眠記録一覧ページ       |      |     x      |    x     |
| [/sleep-records/new](http://localhost:8080/sleep-records/new)          | 睡眠記録入力ページ       |      |     x      |    x     |
| [/sleep-records/{id}](http://localhost:8080/sleep-records/1)           | 睡眠記録詳細ページ       |      |     x      |    x     |
| [/sleep-records/{id}/edit](http://localhost:8080/sleep-records/1/edit) | 睡眠記録編集ページ       |      |     x      |    x     |
| [/statistics](http://localhost:8080/statistics)                        | 統計情報ページ           |      |     x      |    x     |

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

## 開発メモ

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

4. main.goにハンドラーの初期化と登録

    ```go
      // 設定ハンドラーの初期化と登録
      settingsHandler := handler.NewSettingsHandler(tm)
      settingsHandler.RegisterRoutes(router)
    ```

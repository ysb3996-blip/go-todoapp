# go-todoapp

Go の標準ライブラリ(`net/http`・`html/template`・`database/sql`)と SQLite で作った、学習用のシンプルな Todo アプリです。

## できること

- ユーザー登録・ログイン・ログアウト(Cookie とセッション)
- 自分の Todo の一覧・作成・編集・削除
- 他人の Todo は見えない・触れない

## 必要なもの

- Go 1.27 以上(`mise.toml` で 1.27.3 を指定)
- C コンパイラ(SQLite ドライバ `github.com/mattn/go-sqlite3` が cgo を使うため。macOS なら Xcode Command Line Tools)

## 起動

```sh
go run .
```

http://localhost:8080 を開きます。ポート・DB ファイル名・ログファイル名は `config.ini` で変えられます。DB(`webapp.sql`)とテーブルは初回起動時に自動で作られます。

## 構成

```
main.go                 起動
config/                 config.ini の読み込み
utils/                  ログ出力の設定
app/models/             DB 操作(users / todos / sessions)
app/controllers/        URL ごとの処理
app/views/templates/    画面(HTML テンプレート)
app/views/css, js/      Bootstrap 4.5.3・jQuery 3.5.1
```

## 注意

学習用です。パスワードは SHA-1 でハッシュ化しているだけで、本番で使える強度ではありません(本番なら bcrypt などを使う)。CSRF 対策もしていません。

package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	"todoapp/config"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB
var err error

const (
	tableNameUser    = "users"
	tableNameTodo    = "todos"
	tableNameSession = "sessions"
)

func init() {
	Db, err = sql.Open(config.Config.SQLDriver, config.Config.Dbname)
	if err != nil {
		log.Fatalln(err)
	}

	cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid STRING NOT NULL UNIQUE,
		name STRING,
		email STRING,
		password STRING,
		created_at DATETIME)`, tableNameUser)

	_, err = Db.Exec(cmdU)
	if err != nil {
		log.Fatalln(err)
	}

	cmdT := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		user_id INTEGER,
		created_at DATETIME)`, tableNameTodo)

	_, err = Db.Exec(cmdT)
	if err != nil {
		log.Fatalln(err)
	}

	cmdS := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid STRING NOT NULL UNIQUE,
		email STRING,
		user_id INTEGER,
		created_at DATETIME)`, tableNameSession)

	_, err = Db.Exec(cmdS)
	if err != nil {
		log.Fatalln(err)
	}
}

// createUUID は、ユーザーごとに重ならないIDの文字列を作る
func createUUID() string {
	return uuid.New().String()
}

// Encrypt は、パスワードをハッシュ化した文字列にする。
// 注意: SHA-1 は学習用。本番のパスワード保存には bcrypt などを使う
func Encrypt(plaintext string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
}

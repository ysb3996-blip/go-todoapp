package models

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	PassWord  string
	CreatedAt time.Time
}

// Session は sessions テーブルの1行(ログイン中の1人分)
type Session struct {
	ID        int
	UUID      string
	Email     string
	UserID    int
	CreatedAt time.Time
}

// ② ユーザーの作成(Create)
func (u *User) CreateUser() error {
	cmd := `INSERT INTO users (uuid, name, email, password, created_at)
		VALUES (?, ?, ?, ?, ?)`

	_, err := Db.Exec(cmd, createUUID(), u.Name, u.Email, Encrypt(u.PassWord), time.Now())
	return err
}

// ③ ユーザーの取得(Read)
func GetUser(id int) (User, error) {
	user := User{}
	cmd := `SELECT id, uuid, name, email, password, created_at
		FROM users WHERE id = ?`

	err := Db.QueryRow(cmd, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt,
	)
	return user, err
}

// ④ ユーザーの更新(Update)
func (u *User) UpdateUser() error {
	cmd := `UPDATE users SET name = ?, email = ? WHERE id = ?`

	_, err := Db.Exec(cmd, u.Name, u.Email, u.ID)
	return err
}

// ⑤ ユーザーの削除(Delete)
func (u *User) DeleteUser() error {
	cmd := `DELETE FROM users WHERE id = ?`

	_, err := Db.Exec(cmd, u.ID)
	return err
}

// メールアドレスでユーザーを1件取得する(ログイン時に使う)
func GetUserByEmail(email string) (User, error) {
	user := User{}
	cmd := `SELECT id, uuid, name, email, password, created_at
		FROM users WHERE email = ?`

	err := Db.QueryRow(cmd, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt,
	)
	return user, err
}

// ログインしたユーザーのセッションを作って返す
func (u *User) CreateSession() (Session, error) {
	session := Session{}
	sessionUUID := createUUID()

	cmd1 := `INSERT INTO sessions (uuid, email, user_id, created_at)
		VALUES (?, ?, ?, ?)`
	_, err := Db.Exec(cmd1, sessionUUID, u.Email, u.ID, time.Now())
	if err != nil {
		return session, err
	}

	// 今作った1件を uuid で取り直す
	cmd2 := `SELECT id, uuid, email, user_id, created_at
		FROM sessions WHERE uuid = ?`
	err = Db.QueryRow(cmd2, sessionUUID).Scan(
		&session.ID,
		&session.UUID,
		&session.Email,
		&session.UserID,
		&session.CreatedAt,
	)
	return session, err
}

// CheckSession は、そのセッション(uuid)が DB にあるか、つまりログイン中かを調べる。
// 見つかれば sess に DB の内容を詰めて true を返す
func (sess *Session) CheckSession() (bool, error) {
	cmd := `SELECT id, uuid, email, user_id, created_at
		FROM sessions WHERE uuid = ?`

	err := Db.QueryRow(cmd, sess.UUID).Scan(
		&sess.ID,
		&sess.UUID,
		&sess.Email,
		&sess.UserID,
		&sess.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil // セッションが無い = ログインしていない(エラーではない)
	}
	if err != nil {
		return false, err // DB の故障など本当のエラー
	}
	return true, nil
}

// GetUserBySession は、セッションの持ち主のユーザーを返す
func (sess *Session) GetUserBySession() (User, error) {
	user := User{}
	cmd := `SELECT id, uuid, name, email, created_at
		FROM users WHERE id = ?`

	err := Db.QueryRow(cmd, sess.UserID).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	return user, err
}

// DeleteSessionByUUID は、セッションを消す(ログアウト)
func (sess *Session) DeleteSessionByUUID() error {
	cmd := `DELETE FROM sessions WHERE uuid = ?`

	_, err := Db.Exec(cmd, sess.UUID)
	return err
}

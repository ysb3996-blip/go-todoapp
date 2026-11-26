package models

import (
	"time"
)

type Todo struct {
	ID        int
	Content   string
	UserID    int
	CreatedAt time.Time
}

// 作成(Create): ユーザーに紐づく Todo を作る
func (u *User) CreateTodo(content string) error {
	cmd := `INSERT INTO todos (content, user_id, created_at) VALUES (?, ?, ?)`

	_, err := Db.Exec(cmd, content, u.ID, time.Now())
	return err
}

// 取得(Read): id で1件
func GetTodo(id int) (Todo, error) {
	todo := Todo{}
	cmd := `SELECT id, content, user_id, created_at FROM todos WHERE id = ?`

	err := Db.QueryRow(cmd, id).Scan(
		&todo.ID,
		&todo.Content,
		&todo.UserID,
		&todo.CreatedAt,
	)
	return todo, err
}

// 取得(Read): 全件
func GetTodos() ([]Todo, error) {
	cmd := `SELECT id, content, user_id, created_at FROM todos`
	return queryTodos(cmd)
}

// 取得(Read): そのユーザーの Todo だけ
func (u *User) GetTodosByUser() ([]Todo, error) {
	cmd := `SELECT id, content, user_id, created_at FROM todos WHERE user_id = ?`
	return queryTodos(cmd, u.ID)
}

// queryTodos は、複数行の SELECT 結果を []Todo にまとめる
func queryTodos(cmd string, args ...any) ([]Todo, error) {
	rows, err := Db.Query(cmd, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Content,
			&todo.UserID,
			&todo.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, rows.Err()
}

// 更新(Update): 内容と持ち主を更新
func (t *Todo) UpdateTodo() error {
	cmd := `UPDATE todos SET content = ?, user_id = ? WHERE id = ?`

	_, err := Db.Exec(cmd, t.Content, t.UserID, t.ID)
	return err
}

// 削除(Delete)
func (t *Todo) DeleteTodo() error {
	cmd := `DELETE FROM todos WHERE id = ?`

	_, err := Db.Exec(cmd, t.ID)
	return err
}

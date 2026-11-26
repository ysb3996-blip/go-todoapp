package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"todoapp/app/models"
)

func top(w http.ResponseWriter, r *http.Request) {
	// ログイン中なら Todo 一覧へ
	if _, err := session(r); err == nil {
		http.Redirect(w, r, "/todos", http.StatusFound)
		return
	}
	generateHTML(w, nil, "layout", "public_navbar", "top")
}

// index はログイン中のユーザーの Todo 一覧
func index(w http.ResponseWriter, r *http.Request) {
	user, ok := requireUser(w, r)
	if !ok {
		return
	}

	todos, err := user.GetTodosByUser()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		User  models.User
		Todos []models.Todo
	}{user, todos}
	generateHTML(w, data, "layout", "private_navbar", "index")
}

func todoNew(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireUser(w, r); !ok {
		return
	}
	generateHTML(w, nil, "layout", "private_navbar", "todo_new")
}

func todoSave(w http.ResponseWriter, r *http.Request) {
	user, ok := requireUser(w, r)
	if !ok {
		return
	}

	err := user.CreateTodo(r.PostFormValue("content"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/todos", http.StatusFound)
}

func todoEdit(w http.ResponseWriter, r *http.Request) {
	user, ok := requireUser(w, r)
	if !ok {
		return
	}
	todo, ok := ownTodo(w, r, user)
	if !ok {
		return
	}
	generateHTML(w, todo, "layout", "private_navbar", "todo_edit")
}

func todoUpdate(w http.ResponseWriter, r *http.Request) {
	user, ok := requireUser(w, r)
	if !ok {
		return
	}
	todo, ok := ownTodo(w, r, user)
	if !ok {
		return
	}

	todo.Content = r.PostFormValue("content")
	err := todo.UpdateTodo()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/todos", http.StatusFound)
}

func todoDelete(w http.ResponseWriter, r *http.Request) {
	user, ok := requireUser(w, r)
	if !ok {
		return
	}
	todo, ok := ownTodo(w, r, user)
	if !ok {
		return
	}

	err := todo.DeleteTodo()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/todos", http.StatusFound)
}

// ownTodo は URL の {id} の Todo を取り出す。
// 存在しない、または他人の Todo なら 404 を返して ok=false
func ownTodo(w http.ResponseWriter, r *http.Request, user models.User) (todo models.Todo, ok bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return todo, false
	}

	todo, err = models.GetTodo(id)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && todo.UserID != user.ID) {
		http.NotFound(w, r)
		return todo, false
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return todo, false
	}
	return todo, true
}

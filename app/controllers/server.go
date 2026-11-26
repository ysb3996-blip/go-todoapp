package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"todoapp/app/models"
	"todoapp/config"
)

func StartMainServer() error {
	// /static/css/bootstrap.min.css → app/views/css/bootstrap.min.css のように配信する
	files := http.FileServer(http.Dir(config.Config.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("GET /{$}", top)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("GET /login", login)
	http.HandleFunc("POST /authenticate", authenticate)
	http.HandleFunc("GET /logout", logout)

	http.HandleFunc("GET /todos", index)
	http.HandleFunc("GET /todos/new", todoNew)
	http.HandleFunc("POST /todos/save", todoSave)
	http.HandleFunc("GET /todos/edit/{id}", todoEdit)
	http.HandleFunc("POST /todos/update/{id}", todoUpdate)
	http.HandleFunc("POST /todos/delete/{id}", todoDelete)

	return http.ListenAndServe(":"+config.Config.Port, nil)
}

// generateHTML は layout と渡されたテンプレートを組み合わせて画面を返す
func generateHTML(w http.ResponseWriter, data any, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	templates, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Println(err)
	}
}

var errNotLoggedIn = errors.New("not logged in")

// session は Cookie の uuid からセッションを調べる。
// ログインしていなければ errNotLoggedIn を返す
func session(r *http.Request) (models.Session, error) {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		return models.Session{}, errNotLoggedIn
	}

	sess := models.Session{UUID: cookie.Value}
	valid, err := sess.CheckSession()
	if err != nil {
		return sess, err
	}
	if !valid {
		return sess, errNotLoggedIn
	}
	return sess, nil
}

// requireUser はログイン中のユーザーを返す。
// ログインしていなければ /login に飛ばして ok=false を返す
func requireUser(w http.ResponseWriter, r *http.Request) (user models.User, ok bool) {
	sess, err := session(r)
	if errors.Is(err, errNotLoggedIn) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return user, false
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return user, false
	}

	user, err = sess.GetUserBySession()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return user, false
	}
	return user, true
}

package controllers

import (
	"log"
	"net/http"
	"todoapp/app/models"
)

func signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 登録フォームを表示する
		generateHTML(w, nil, "layout", "public_navbar", "signup")

	case http.MethodPost:
		// フォームの入力からユーザーを作る
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		user := models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			PassWord: r.PostFormValue("password"),
		}
		err = user.CreateUser()
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	// ログイン中なら Todo 一覧へ
	if _, err := session(r); err == nil {
		http.Redirect(w, r, "/todos", http.StatusFound)
		return
	}
	generateHTML(w, nil, "layout", "public_navbar", "login")
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess := models.Session{UUID: cookie.Value}
		err = sess.DeleteSessionByUUID()
		if err != nil {
			log.Println(err)
		}
	}

	// ブラウザの Cookie も消す
	http.SetCookie(w, &http.Cookie{
		Name:     "_cookie",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// メールアドレスでユーザーを探す。いなければログイン画面に戻す
	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// パスワードが違えばログイン画面に戻す
	if user.PassWord != models.Encrypt(r.PostFormValue("password")) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// セッションを作って、その uuid を Cookie に入れる
	session, err := user.CreateSession()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    session.UUID,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusFound)
}

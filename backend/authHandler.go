package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	// "net/http"
)

// +++++++++++++++++structs+++++++++++++++++
type Info struct {
	Username string `json:username`
	Password string `json:password`
	Email    string `json:email`
}

//+++++++++++++++++structs++++++++++++++++++

// Global session store
var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

// ++++++++++++++ New Helper Function ++++++++++++++
// This function applies the session and is called by both register and login
func applySession(w http.ResponseWriter, req *http.Request, user *Info) {

	session, err := store.Get(req, "authentication")

	fmt.Println(session.IsNew)

	if err != nil && session.IsNew != true {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	session.Options = &sessions.Options{HttpOnly: false, SameSite: http.SameSiteDefaultMode, Secure: false, MaxAge: 2628000, Path: "/"}

	session.ID = fmt.Sprintf("%x", md5.Sum([]byte(user.Email+time.Now().String()+string(securecookie.GenerateRandomKey(16)))))
	session.Values["Authenticated"] = true

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO sessions (user_id, session_token, expires_at) VALUES ((SELECT id FROM users WHERE email = ?), ?, DATE_ADD(NOW(), INTERVAL 1 MONTH)) ON DUPLICATE KEY UPDATE expires_at = DATE_ADD(NOW(), INTERVAL 1 MONTH)")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(user.Email, session.ID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := session.Save(req, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
}

// ++++++++++++++ Register Handler ++++++++++++++
func register(w http.ResponseWriter, req *http.Request) {
	var userInfo Info
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if userInfo.Username == "" || userInfo.Email == "" || userInfo.Password == "" {
		http.Error(w, "Inputs can't be empty", http.StatusBadRequest)
		return
	}
	passwordHash := fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Password)))
	// db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	// if err != nil {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }
	// defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)")
	// if err != nil {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }
	// defer stmt.Close()
	if addValues("users", []string{"username", "email", "password_hash"}, [][]string{{userInfo.Username, userInfo.Email, passwordHash}}) != nil {
		http.Error(w, "Username or email already exists", http.StatusConflict)
		return
	}

	// _, err = stmt.Exec(userInfo.Username, userInfo.Email, passwordHash)
	// if err != nil {
	// 	http.Error(w, "Username or email already exists", http.StatusConflict)
	// 	return
	// }

	applySession(w, req, &userInfo)

	w.WriteHeader(http.StatusCreated)
}

// ++++++++++++++ Login Handler ++++++++++++++
func login(w http.ResponseWriter, req *http.Request) {

	var userInfo Info
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	row := db.QueryRow("SELECT username, email FROM users WHERE email = ? AND password_hash = ?", userInfo.Email, fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Password))))

	var storedUsername, storedEmail string
	if err := row.Scan(&storedUsername, &storedEmail); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	applySession(w, req, &Info{Username: storedUsername, Email: storedEmail})

	w.WriteHeader(http.StatusOK)
}

func getCookieUser(w http.ResponseWriter, req *http.Request) {

	cookie, err := req.Cookie("authentication")
	if err != nil {
		fmt.Println("Error getting cookie:", err)
		return
	}
	if cookie.Value == "" {
		fmt.Println("Cookie value is empty")
		return
	}
	// cookieValue := cookie.Value
	var session *sessions.Session
	if cookie.Value != "" {

		session, err = store.Get(req, "authentication")
		if err != nil {
			fmt.Println("Error getting session:", err)
			return
		}

		auth, ok := session.Values["Authenticated"].(bool)
		if !ok || !auth {
			fmt.Println("Unauthorized")
			return
		}

		user, err := findValues("users", []string{"username", "email"}, map[string]string{"email": session.Values["email"].(string)})
		if err != nil {
			fmt.Println("Internal server error:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user[0])
	}
}

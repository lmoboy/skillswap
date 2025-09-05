package onlineChat

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
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userInfo.Username, userInfo.Email, passwordHash)
	if err != nil {
		http.Error(w, "Username or email already exists", http.StatusConflict)
		return
	}

	// Call the helper function to set the session
	applySession(w, req, &userInfo)

	// session, err := store.Get(req, "fuckasssession")
	// if err != nil {
	// 	http.Error(w, "Failed to get session", http.StatusInternalServerError)
	// 	return
	// }

	// session.Values["token"] = "session.id"
	// session.Values["Authenticated"] = true
	// session.Values["username"] = userInfo.Username
	// session.Values["email"] = userInfo.Email

	// if err := session.Save(req, w); err != nil {
	// 	http.Error(w, "Failed to save session", http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
	// fmt.Fprintln(w, "Registration successful. Session set.")
}

// ++++++++++++++ Login Handler ++++++++++++++
func login(w http.ResponseWriter, req *http.Request) {

	var userInfo Info
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Use prepared statements for safe querying
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	// fmt.Println(userInfo)
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
	// session, err := store.Get(req, "fuckasssession")
	// if err != nil {
	// 	http.Error(w, "Failed to get session", http.StatusInternalServerError)
	// 	return
	// }

	// session.Values["token"] = "session.id"
	// session.Values["Authenticated"] = true
	// session.Values["username"] = userInfo.Username
	// session.Values["email"] = userInfo.Email

	// if err := session.Save(req, w); err != nil {
	// 	http.Error(w, "Failed to save session", http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	// fmt.Fprintln(w, "Login successful. Session set.")
}

func checkValidAuth(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "authentication")
	if err != nil {
		http.Error(w, "no session found", http.StatusInternalServerError)
	}

	auth, ok := session.Values["Authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User is authenticated")
}

func main() {

}

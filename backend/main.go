package main

import (
	"net/http"
	"time"

	"skillswap/backend/authentication"
	"skillswap/backend/chat"
	"skillswap/backend/config"
	"skillswap/backend/database"
	"skillswap/backend/users"
	"skillswap/backend/utils"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

// main initializes the application (database, routes and middleware) and starts the HTTP server.
// It registers API endpoints for authentication, user, chat and search operations, configures CORS, and listens on localhost:8080 with 15s read and write timeouts.
func main() {
	database.Init()
	// Izveido jaunu rūteri ar stingru pārbaudi slīpsvītrām, kas nozīmē, ka maršruti ar un bez beigu slīpsvītras tiek uzskatīti par atšķirīgiem.
	server := mux.NewRouter().StrictSlash(true)

	// Tiek definēti API ceļi (end-points) dažādām front-end darbībām.
	// "HandleFunc" piesaista konkrētu URL ceļu noteiktai Go funkcijai.
	// server.HandleFunc("/api/chat", RunWebsocket)
	server.HandleFunc("/api/login", authentication.Login).Methods("POST")
	server.HandleFunc("/api/register", authentication.Register).Methods("POST")
	server.HandleFunc("/api/logout", authentication.Logout).Methods("POST")
	server.HandleFunc("/api/cookieUser", authentication.CheckSession).Methods("GET")
	server.HandleFunc("/api/search", database.Search).Methods("POST")
	server.HandleFunc("/api/user", users.RetrieveUserInfo).Methods("GET")

	server.HandleFunc("/api/getChats", chat.GetChatsFromUserID)
	server.HandleFunc("/api/getChatInfo", chat.GetMessagesFromUID)
	// server.HandleFunc("/api/video", websocket.JoinWebSocket).Methods("GET")

	// Vienkārša "dummy" funkcija aizmugursistēmas (backend) darbības pārbaudei.
	// Tā atgriež JSON atbildi ar statusu "pong", kad tiek saņemts GET pieprasījums.
	server.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "pong"})
	}).Methods("GET")

	// Tiek konfigurēts CORS (Cross-Origin Resource Sharing), lai atļautu pieprasījumus
	// no noteiktiem "oriģiniem" (front-end servera adresēm).
	c := config.CORS()

	// Piesaista CORS apstrādātāju rūterim, lai tas darbotos ar visiem ceļiem.
	server.Use(c.Handler)

	// Reģistrē rūteri kā galveno HTTP apstrādātāju.
	http.Handle("/", server)

	// Izveido un konfigurē HTTP serveri.
	serve := &http.Server{
		Handler: server,           // Norāda, ka šis serveris izmantos mūsu iepriekš definēto rūteri.
		Addr:    "localhost:8080", // Serveris klausīsies uz šīs adreses un porta.

		// Iestata maksimālo laiku, cik ilgi serveris gaidīs pieprasījuma ķermeņa nolasīšanu un atbildes uzrakstīšanu.
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Sāk servera darbību. "ListenAndServe" bloķē izpildi, līdz serveris tiek apturēts.
	serve.ListenAndServe()
}

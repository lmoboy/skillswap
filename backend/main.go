package main

import (
	"net/http"
	"time"

	"skillswap/backend/authentication"
	"skillswap/backend/chat"
	"skillswap/backend/config"
	"skillswap/backend/courses"
	"skillswap/backend/database"
	"skillswap/backend/structs"
	"skillswap/backend/users"
	"skillswap/backend/utils"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

func getSkills(w http.ResponseWriter, req *http.Request) {
	rows, err := database.Query(`SELECT id,name,description FROM skills`)
	if err != nil {
		utils.HandleError(err)
		return
	}
	defer rows.Close()
	skills := []structs.Skill{}
	for rows.Next() {
		var p structs.Skill
		rows.Scan(&p.ID, &p.Name, &p.Description)
		skills = append(skills, p)
	}

	utils.SendJSONResponse(w, http.StatusOK, skills)
}

// main initializes the application (database, routes and middleware) and starts the HTTP server.
// main initializes the database, registers API endpoints for authentication, user, chat, and search, configures CORS, and starts an HTTP server on localhost:8080.
// The server uses the configured router with ReadTimeout and WriteTimeout set to 15 seconds and blocks while serving.
// main initializes the database, starts the chat WebSocket hub, registers HTTP API routes (authentication, user, chat, search, and health), configures CORS, and starts an HTTP server on localhost:8080 with 15-second read and write timeouts.
// Registered endpoints handle login/register/logout, session check, search and fullSearch, user info retrieval, websocket and chat operations, and a /api/ping health check.
func main() {
	database.Init()

	// Start the WebSocket hub for chat functionality
	go chat.StartHub()

	// Izveido jaunu rūteri ar stingru pārbaudi slīpsvītrām, kas nozīmē, ka maršruti ar un bez beigu slīpsvītras tiek uzskatīti par atšķirīgiem.
	server := mux.NewRouter().StrictSlash(true)

	// Tiek definēti API ceļi (end-points) dažādām front-end darbībām.
	// "HandleFunc" piesaista konkrētu URL ceļu noteiktai Go funkcijai.
	// server.HandleFunc("/api/chat", chat.RunWebsocket)
	server.HandleFunc("/api/login", authentication.Login).Methods("POST")
	server.HandleFunc("/api/register", authentication.Register).Methods("POST")
	server.HandleFunc("/api/logout", authentication.Logout).Methods("POST")
	server.HandleFunc("/api/cookieUser", authentication.CheckSession).Methods("GET")
	server.HandleFunc("/api/updateUser", users.UpdateUser).Methods("POST")
	server.HandleFunc("/api/profile/picture", users.UploadProfilePicture).Methods("POST")
	server.HandleFunc("/api/profile/{id}/picture", users.GetProfilePicture).Methods("GET")

	server.HandleFunc("/api/search", database.Search).Methods("POST")
	server.HandleFunc("/api/fullSearch", database.FullSearch).Methods("POST")
	server.HandleFunc("/api/user", users.RetrieveUserInfo).Methods("GET")

	server.HandleFunc("/api/sendMessage", chat.SaveToDBLink).Methods("POST")

	server.HandleFunc("/api/chat", chat.SimpleWebSocketEndpoint)
	server.HandleFunc("/api/createChat", chat.CreateChat)
	server.HandleFunc("/api/getChats", chat.GetChatsFromUserID)
	server.HandleFunc("/api/getChatInfo", chat.GetMessagesFromUID)
	// server.HandleFunc("/api/video", websocket.JoinWebSocket).Methods("GET")

	server.HandleFunc("/api/courses", courses.GetAllCourses).Methods("GET")
	server.HandleFunc("/api/course", courses.GetCourseByID).Methods("GET")
	server.HandleFunc("/api/course/upload", courses.UploadCourseAsset).Methods("POST")
	server.HandleFunc("/api/course/{id}/stream", courses.StreamCourseAsset).Methods("GET")
	server.HandleFunc("/api/searchCourses", courses.SearchCourses).Methods("POST")
	server.HandleFunc("/api/coursesByInstructor", courses.GetCoursesByInstructor).Methods("GET")

	server.HandleFunc("/api/getSkills", getSkills)

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

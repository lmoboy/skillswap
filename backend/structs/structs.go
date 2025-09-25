package structs

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type UserInfo struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ID             string `json:"id"`
	ProfilePicture string `json:"profile_picture"`
	AboutMe        string `json:"aboutme"`
	Projects       string `json:"projects"`
	Professions    string `json:"professions"`
	Contacts       string `json:"contacts"`
	Skills         string `json:"skills"`
	Location       string `json:"location"`
}
type Skill struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type SessionData struct {
	UserInfo *UserInfo
	Session  *sessions.Session
	Cookie   *http.Cookie
}

type SearchResult struct {
	User        UserInfo `json:"user"`
	SkillsFound string   `json:"skills_found"`
}

package structs

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type UserInfo struct {
	Username       string        `json:"username"`
	Email          string        `json:"email"`
	Password       string        `json:"password"`
	ID             int           `json:"id"`
	ProfilePicture string        `json:"profile_picture"`
	AboutMe        string        `json:"aboutme"`
	Projects       []UserProject `json:"projects"`
	Professions    string        `json:"profession"`
	Contacts       []UserContact `json:"contacts"`
	Skills         []UserSkill   `json:"skills"`
	Location       string        `json:"location"`
	Joined         string        `json:"created_at"`
}

type UserProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

type UserSkill struct {
	Name     string `json:"name"`
	Verified int   `json:"verified"`
}

type UserContact struct {
	Name string `json:"name"`
	Link string `json:"link"`
	Icon string `json:"icon"`
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

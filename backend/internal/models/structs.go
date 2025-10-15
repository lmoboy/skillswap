package models

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

// FlexBool is a custom type that can unmarshal both boolean and integer values
type FlexBool bool

// UnmarshalJSON implements custom unmarshaling for FlexBool
func (fb *FlexBool) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as boolean first
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		*fb = FlexBool(b)
		return nil
	}

	// Try to unmarshal as integer
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*fb = FlexBool(i != 0)
		return nil
	}

	// Try to unmarshal as string (in case it comes as "0", "1", "true", "false")
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if b, err := strconv.ParseBool(s); err == nil {
			*fb = FlexBool(b)
			return nil
		}
		if i, err := strconv.Atoi(s); err == nil {
			*fb = FlexBool(i != 0)
			return nil
		}
	}

	// Default to false if all parsing fails
	*fb = false
	return nil
}

// MarshalJSON implements custom marshaling for FlexBool
func (fb FlexBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(fb))
}

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
	Name     string   `json:"name"`
	Verified FlexBool `json:"verified"`
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

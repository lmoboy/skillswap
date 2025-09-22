package structs

type UserInfo struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ID             int    `json:"id"`
	ProfilePicture string `json:"profile_picture"`
	AboutMe        string `json:"aboutme"`
	Projects       string `json:"projects"`
	Contacts       string `json:"contacts"`
	Skills         string `json:"skills"`
	Location       string `json:"location"`
}
type Skill struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
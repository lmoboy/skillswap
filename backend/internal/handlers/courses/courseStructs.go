package courses

type Course struct {
	ID              int64   `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	InstructorID    int64   `json:"instructor_id"`
	InstructorName  string  `json:"instructor_name"`
	SkillID         int64   `json:"skill_id"`
	SkillName       string  `json:"skill_name"`
	DifficultyLevel string  `json:"difficulty_level"`
	DurationHours   int     `json:"duration_hours"`
	MaxStudents     int     `json:"max_students"`
	CurrentStudents int     `json:"current_students"`
	Price           float64 `json:"price"`
	ThumbnailURL    string  `json:"thumbnail_url"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	AverageRating   float64 `json:"average_rating"`
	ReviewCount     int     `json:"review_count"`
}

type CourseModule struct {
	ID            int64  `json:"id"`
	CourseID      int64  `json:"course_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	VideoURL      string `json:"video_url"`
	VideoDuration int    `json:"video_duration"`
	ThumbnailURL  string `json:"thumbnail_url"`
	OrderIndex    int    `json:"order_index"`
	CreatedAt     string `json:"created_at"`
}

type CourseEnrollment struct {
	ID          int64  `json:"id"`
	CourseID    int64  `json:"course_id"`
	StudentID   int64  `json:"student_id"`
	EnrolledAt  string `json:"enrolled_at"`
	CompletedAt string `json:"completed_at"`
	Progress    int    `json:"progress"`
}

type CourseReview struct {
	ID          int64  `json:"id"`
	CourseID    int64  `json:"course_id"`
	StudentID   int64  `json:"student_id"`
	StudentName string `json:"student_name"`
	Rating      int    `json:"rating"`
	ReviewText  string `json:"review_text"`
	CreatedAt   string `json:"created_at"`
}

type CourseDetail struct {
	Course
	Modules []CourseModule `json:"modules"`
	Reviews []CourseReview `json:"reviews"`
}

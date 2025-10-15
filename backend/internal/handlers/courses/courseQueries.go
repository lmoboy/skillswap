package courses

import (
	"database/sql"
	"skillswap/backend/internal/database"
)

const (
	// Course query with ratings
	courseWithRatingsQuery = `
		SELECT
			c.id, c.title, c.description, c.instructor_id, u.username as instructor_name,
			c.skill_id, s.name as skill_name, c.difficulty_level, c.duration_hours,
			c.max_students, c.current_students, c.price, c.thumbnail_url, c.status,
			c.created_at, COALESCE(c.updated_at, ''),
			COALESCE(AVG(cr.rating), 0) as average_rating,
			COUNT(cr.id) as review_count
		FROM courses c
		JOIN users u ON c.instructor_id = u.id
		JOIN skills s ON c.skill_id = s.id
		LEFT JOIN course_reviews cr ON c.id = cr.course_id
	`

	courseModulesQuery = `
		SELECT id, course_id, title, description, COALESCE(video_url, ''), 
		       video_duration, COALESCE(thumbnail_url, ''), order_index, created_at
		FROM course_modules
		WHERE course_id = ?
		ORDER BY order_index ASC
	`

	courseReviewsQuery = `
		SELECT cr.id, cr.course_id, cr.student_id, u.username as student_name, 
		       cr.rating, cr.review_text, cr.created_at
		FROM course_reviews cr
		JOIN users u ON cr.student_id = u.id
		WHERE cr.course_id = ?
		ORDER BY cr.created_at DESC
		LIMIT 10
	`
)

// scanCourse scans a course row into a Course struct
func scanCourse(scanner interface {
	Scan(dest ...interface{}) error
}, course *Course) error {
	return scanner.Scan(
		&course.ID, &course.Title, &course.Description, &course.InstructorID, &course.InstructorName,
		&course.SkillID, &course.SkillName, &course.DifficultyLevel, &course.DurationHours,
		&course.MaxStudents, &course.CurrentStudents, &course.Price, &course.ThumbnailURL, &course.Status,
		&course.CreatedAt, &course.UpdatedAt, &course.AverageRating, &course.ReviewCount,
	)
}

// fetchCourses retrieves multiple courses based on a query and args
func fetchCourses(query string, args ...interface{}) ([]Course, error) {
	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := scanCourse(rows, &course); err == nil {
			courses = append(courses, course)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

// fetchCourseByID retrieves a single course by ID
func fetchCourseByID(courseID int64) (*Course, error) {
	var course Course
	query := courseWithRatingsQuery + `
		WHERE c.id = ? AND c.status = 'Published'
		GROUP BY c.id
	`
	
	err := database.QueryRow(query, courseID).Scan(
		&course.ID, &course.Title, &course.Description, &course.InstructorID, &course.InstructorName,
		&course.SkillID, &course.SkillName, &course.DifficultyLevel, &course.DurationHours,
		&course.MaxStudents, &course.CurrentStudents, &course.Price, &course.ThumbnailURL, &course.Status,
		&course.CreatedAt, &course.UpdatedAt, &course.AverageRating, &course.ReviewCount,
	)

	if err != nil {
		return nil, err
	}
	return &course, nil
}

// fetchCourseModules retrieves all modules for a course
func fetchCourseModules(courseID int64) ([]CourseModule, error) {
	rows, err := database.Query(courseModulesQuery, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []CourseModule
	for rows.Next() {
		var module CourseModule
		err := rows.Scan(&module.ID, &module.CourseID, &module.Title, &module.Description,
			&module.VideoURL, &module.VideoDuration, &module.ThumbnailURL,
			&module.OrderIndex, &module.CreatedAt)
		if err == nil {
			modules = append(modules, module)
		}
	}

	return modules, nil
}

// fetchCourseReviews retrieves reviews for a course
func fetchCourseReviews(courseID int64) ([]CourseReview, error) {
	rows, err := database.Query(courseReviewsQuery, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []CourseReview
	for rows.Next() {
		var review CourseReview
		err := rows.Scan(&review.ID, &review.CourseID, &review.StudentID, &review.StudentName,
			&review.Rating, &review.ReviewText, &review.CreatedAt)
		if err == nil {
			reviews = append(reviews, review)
		}
	}

	return reviews, nil
}

// getSkillIDByName retrieves the skill ID for a given skill name
func getSkillIDByName(skillName string) (int64, error) {
	var skillID int64
	err := database.QueryRow("SELECT id FROM skills WHERE name = ?", skillName).Scan(&skillID)
	if err == sql.ErrNoRows {
		return 0, sql.ErrNoRows
	}
	return skillID, err
}

// insertCourse creates a new course in the database
func insertCourse(title, description string, instructorID, skillID int64, durationHours int, thumbnailURL string) (int64, error) {
	res, err := database.Execute(`
		INSERT INTO courses (title, description, instructor_id, skill_id, duration_hours, thumbnail_url, status)
		VALUES (?, ?, ?, ?, ?, ?, 'Published')
	`, title, description, instructorID, skillID, durationHours, thumbnailURL)
	
	if err != nil {
		return 0, err
	}
	
	return res.LastInsertId()
}

// insertCourseModule creates a new course module in the database
func insertCourseModule(courseID int64, title, description, videoURL string, videoDuration int, thumbnailURL string, orderIndex int) error {
	_, err := database.Execute(`
		INSERT INTO course_modules (course_id, title, description, video_url, video_duration, thumbnail_url, order_index)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, courseID, title, description, videoURL, videoDuration, thumbnailURL, orderIndex)
	
	return err
}



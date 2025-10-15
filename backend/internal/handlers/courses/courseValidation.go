package courses

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrInvalidDuration       = errors.New("invalid duration value")
	ErrInvalidID             = errors.New("invalid ID")
)

// CourseFormData represents the basic course information from the form
type CourseFormData struct {
	Title           string
	Description     string
	SkillName       string
	DurationMinutes int
}

// ModuleFormData represents a single module's data
type ModuleFormData struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	OrderIndex    int    `json:"order_index"`
	VideoDuration int    `json:"video_duration"`
}

// validateCourseFormData validates basic course form fields
func validateCourseFormData(r *http.Request) (*CourseFormData, error) {
	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))
	skillName := strings.TrimSpace(r.FormValue("skill_name"))
	durationMinutesStr := r.FormValue("duration_minutes")

	if title == "" || description == "" || skillName == "" || durationMinutesStr == "" {
		return nil, ErrMissingRequiredFields
	}

	durationMinutes, err := strconv.Atoi(durationMinutesStr)
	if err != nil {
		return nil, ErrInvalidDuration
	}

	return &CourseFormData{
		Title:           title,
		Description:     description,
		SkillName:       skillName,
		DurationMinutes: durationMinutes,
	}, nil
}

// parseModulesFromForm parses module data from the form JSON
func parseModulesFromForm(r *http.Request) ([]ModuleFormData, error) {
	modulesJSON := r.FormValue("modules")
	if modulesJSON == "" {
		return []ModuleFormData{}, nil
	}

	var modules []ModuleFormData
	err := json.Unmarshal([]byte(modulesJSON), &modules)
	if err != nil {
		return nil, err
	}

	return modules, nil
}

// validateAndParseInt64 validates and parses an int64 from a string
func validateAndParseInt64(value string) (int64, error) {
	if value == "" {
		return 0, ErrInvalidID
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, ErrInvalidID
	}

	return id, nil
}



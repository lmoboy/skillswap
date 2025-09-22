package database

// GetUserIDFromEmail returns the user ID for a given email
func GetUserIDFromEmail(email string) (int64, error) {
	db, err := getDatabase()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	var id int64
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// GetSkillIDFromName returns the skill ID for a given name
func GetSkillIDFromName(name string) (int64, error) {
	db, err := getDatabase()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	var id int64
	err = db.QueryRow("SELECT id FROM skills WHERE name = ?", name).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// GetAllSkills returns a list of all skills in the database
func GetAllSkills() ([]Skill, error) {
	db, err := getDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, description FROM skills")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []Skill
	for rows.Next() {
		var skill Skill
		err := rows.Scan(&skill.ID, &skill.Name, &skill.Description)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return skills, nil
}

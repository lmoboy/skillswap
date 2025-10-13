package structs

import (
	"encoding/json"
	"testing"
)

func TestFlexBoolUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		wantErr  bool
	}{
		// Boolean inputs
		{
			name:     "boolean true",
			input:    `{"verified": true}`,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "boolean false",
			input:    `{"verified": false}`,
			expected: false,
			wantErr:  false,
		},
		// Integer inputs (MySQL TINYINT)
		{
			name:     "integer 1",
			input:    `{"verified": 1}`,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "integer 0",
			input:    `{"verified": 0}`,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "integer 42",
			input:    `{"verified": 42}`,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "negative integer",
			input:    `{"verified": -1}`,
			expected: true,
			wantErr:  false,
		},
		// String inputs
		{
			name:     "string true",
			input:    `{"verified": "true"}`,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "string false",
			input:    `{"verified": "false"}`,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "string 1",
			input:    `{"verified": "1"}`,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "string 0",
			input:    `{"verified": "0"}`,
			expected: false,
			wantErr:  false,
		},
		// Edge cases
		{
			name:     "null value defaults to false",
			input:    `{"verified": null}`,
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result struct {
				Verified FlexBool `json:"verified"`
			}

			err := json.Unmarshal([]byte(tt.input), &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if bool(result.Verified) != tt.expected {
				t.Errorf("UnmarshalJSON() got = %v, want %v", bool(result.Verified), tt.expected)
			}
		})
	}
}

func TestFlexBoolMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    FlexBool
		expected string
	}{
		{
			name:     "true value",
			input:    FlexBool(true),
			expected: "true",
		},
		{
			name:     "false value",
			input:    FlexBool(false),
			expected: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}

			if string(result) != tt.expected {
				t.Errorf("MarshalJSON() got = %s, want %s", string(result), tt.expected)
			}
		})
	}
}

func TestUserSkillWithFlexBool(t *testing.T) {
	// Test full UserSkill unmarshaling with different verified formats
	tests := []struct {
		name     string
		input    string
		expected UserSkill
	}{
		{
			name:  "skill with boolean verified",
			input: `{"name": "JavaScript", "verified": true}`,
			expected: UserSkill{
				Name:     "JavaScript",
				Verified: FlexBool(true),
			},
		},
		{
			name:  "skill with integer verified",
			input: `{"name": "Python", "verified": 1}`,
			expected: UserSkill{
				Name:     "Python",
				Verified: FlexBool(true),
			},
		},
		{
			name:  "skill with zero verified",
			input: `{"name": "Go", "verified": 0}`,
			expected: UserSkill{
				Name:     "Go",
				Verified: FlexBool(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var skill UserSkill
			err := json.Unmarshal([]byte(tt.input), &skill)
			if err != nil {
				t.Errorf("Unmarshal UserSkill error = %v", err)
				return
			}

			if skill.Name != tt.expected.Name {
				t.Errorf("Name got = %s, want %s", skill.Name, tt.expected.Name)
			}

			if bool(skill.Verified) != bool(tt.expected.Verified) {
				t.Errorf("Verified got = %v, want %v", bool(skill.Verified), bool(tt.expected.Verified))
			}
		})
	}
}

func TestUserSkillArrayUnmarshal(t *testing.T) {
	// Test unmarshaling an array of skills with mixed verified types
	input := `[
		{"name": "JavaScript", "verified": true},
		{"name": "Python", "verified": 1},
		{"name": "Go", "verified": 0},
		{"name": "Rust", "verified": false}
	]`

	var skills []UserSkill
	err := json.Unmarshal([]byte(input), &skills)
	if err != nil {
		t.Fatalf("Unmarshal skills array error = %v", err)
	}

	if len(skills) != 4 {
		t.Errorf("Expected 4 skills, got %d", len(skills))
	}

	expectedResults := []bool{true, true, false, false}
	for i, skill := range skills {
		if bool(skill.Verified) != expectedResults[i] {
			t.Errorf("Skill %d (%s) verified got = %v, want %v",
				i, skill.Name, bool(skill.Verified), expectedResults[i])
		}
	}
}

func TestFlexBoolComparison(t *testing.T) {
	// Test that FlexBool can be compared with bool
	trueFlexBool := FlexBool(true)
	falseFlexBool := FlexBool(false)

	if bool(trueFlexBool) != true {
		t.Error("FlexBool(true) should equal true")
	}

	if bool(falseFlexBool) != false {
		t.Error("FlexBool(false) should equal false")
	}

	// Test comparison in conditional
	if bool(trueFlexBool) {
		// Expected path
	} else {
		t.Error("FlexBool(true) should be truthy")
	}

	if bool(falseFlexBool) {
		t.Error("FlexBool(false) should be falsy")
	}
}

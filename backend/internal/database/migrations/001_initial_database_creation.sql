-- Initial database schema for SkillSwap
-- This migration creates all core tables
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS course_reviews;
DROP TABLE IF EXISTS course_enrollments;
DROP TABLE IF EXISTS course_modules;
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS user_contacts;
DROP TABLE IF EXISTS user_projects;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS user_skills;
DROP TABLE IF EXISTS skills;
DROP TABLE IF EXISTS users;

SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  username VARCHAR(50) NOT NULL,
  email VARCHAR(191) NOT NULL,
  profile_picture VARCHAR(255) DEFAULT "noPicture",
  aboutme TEXT,
  profession TEXT,
  location VARCHAR(191) DEFAULT "",
  swaps INT NOT NULL DEFAULT 2,

  password_hash VARCHAR(255) NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  UNIQUE KEY uq_users_email (email),
  KEY idx_users_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS skills (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_skills_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_skills (
  user_id BIGINT UNSIGNED NOT NULL,
  skill_id BIGINT UNSIGNED NOT NULL,
  teaching_skill ENUM(
    'Show how it''s done',
    'Explain in slight details',
    'Give homework report on subject',
    'Professor'
  ) NOT NULL DEFAULT 'Show how it''s done',
  verified BOOL NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (user_id, skill_id),
  CONSTRAINT fk_user_skills_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_skills_skill
    FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE,

  KEY idx_user_skills_skill (skill_id),
  KEY idx_user_skills_user_verified (user_id, verified)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS chats (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user1_id BIGINT UNSIGNED NOT NULL,
  user2_id BIGINT UNSIGNED NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  CONSTRAINT fk_chats_user1 FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_chats_user2 FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE,

  UNIQUE KEY uq_chats_pair (user1_id, user2_id),

  KEY idx_chats_user1 (user1_id),
  KEY idx_chats_user2 (user2_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS messages (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  chat_id BIGINT UNSIGNED NOT NULL,
  sender_id BIGINT UNSIGNED NOT NULL,
  content TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  CONSTRAINT fk_messages_chat FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
  CONSTRAINT fk_messages_sender FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,

  KEY idx_messages_chat (chat_id),
  KEY idx_messages_sender_time (sender_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_projects (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(191) NOT NULL,
  description TEXT,
  link VARCHAR(255) DEFAULT "",
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  CONSTRAINT fk_user_projects_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

  KEY idx_user_projects_user (user_id),
  KEY idx_user_projects_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_contacts (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(191) NOT NULL,
  link VARCHAR(255) DEFAULT "",
  icon VARCHAR(255) NOT NULL DEFAULT("MessageCircleQuestionMark"),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  CONSTRAINT fk_user_contacts_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

  KEY idx_user_contacts_user (user_id),
  KEY idx_user_contacts_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS courses (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  instructor_id BIGINT UNSIGNED NOT NULL,
  skill_id BIGINT UNSIGNED NOT NULL,
  difficulty_level ENUM('Beginner', 'Intermediate', 'Advanced', 'Expert') NOT NULL DEFAULT 'Beginner',
  duration_hours INT UNSIGNED DEFAULT 0,
  max_students INT UNSIGNED DEFAULT 10,
  current_students INT UNSIGNED DEFAULT 0,
  price DECIMAL(10, 2) DEFAULT 0.00,
  thumbnail_url VARCHAR(255) DEFAULT "",
  status ENUM('Draft', 'Published', 'Archived') NOT NULL DEFAULT 'Draft',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  CONSTRAINT fk_courses_instructor
    FOREIGN KEY (instructor_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_courses_skill
    FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE,

  KEY idx_courses_instructor (instructor_id),
  KEY idx_courses_skill (skill_id),
  KEY idx_courses_status (status),
  KEY idx_courses_difficulty (difficulty_level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS course_modules (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  course_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  video_url VARCHAR(500),
  video_duration INT UNSIGNED DEFAULT 0,
  thumbnail_url VARCHAR(500),
  order_index INT UNSIGNED NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  CONSTRAINT fk_course_modules_course
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,

  KEY idx_course_modules_course (course_id),
  KEY idx_course_modules_order (course_id, order_index)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS course_enrollments (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  course_id BIGINT UNSIGNED NOT NULL,
  student_id BIGINT UNSIGNED NOT NULL,
  enrolled_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  completed_at TIMESTAMP NULL DEFAULT NULL,
  progress INT UNSIGNED DEFAULT 0,

  PRIMARY KEY (id),
  CONSTRAINT fk_course_enrollments_course
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
  CONSTRAINT fk_course_enrollments_student
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,

  UNIQUE KEY uq_enrollment (course_id, student_id),
  KEY idx_course_enrollments_student (student_id),
  KEY idx_course_enrollments_course (course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS course_reviews (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  course_id BIGINT UNSIGNED NOT NULL,
  student_id BIGINT UNSIGNED NOT NULL,
  rating INT UNSIGNED NOT NULL CHECK (rating >= 1 AND rating <= 5),
  review_text TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  CONSTRAINT fk_course_reviews_course
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
  CONSTRAINT fk_course_reviews_student
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,

  UNIQUE KEY uq_review (course_id, student_id),
  KEY idx_course_reviews_course (course_id),
  KEY idx_course_reviews_rating (rating)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Initial seed data for skills
INSERT IGNORE INTO skills (name, description) VALUES
    ('Python', 'Programming language for general-purpose development'),
    ('JavaScript', 'Language for web development'),
    ('SQL', 'Database query language'),
    ('Graphic Design', 'Creating visual content'),
    ('Public Speaking', 'Skill for effective oral communication'),
    ('React', 'JavaScript library for building user interfaces'),
    ('Svelte', 'JavaScript framework for building web applications'),
    ('Tailwind CSS', 'Utility-first CSS framework for rapid development'),
    ('Node.js', 'JavaScript runtime for developing server-side applications'),
    ('C++', 'Programming language for developing operating systems and applications'),
    ('Java', 'Programming language for developing enterprise applications'),
    ('HTML/CSS', 'Markup language and styling language for web development'),
    ('Git', 'Version control system for tracking changes in source code'),
    ('Rust', 'Programming language for developing system software'),
    ('Kotlin', 'Programming language for developing Android applications'),
    ('TypeScript', 'Statically typed JavaScript superset'),
    ('Machine Learning', 'Development of algorithms for prediction and classification'),
    ('Computer Vision', 'Development of algorithms for image and video analysis'),
    ('Natural Language Processing', 'Development of algorithms for text analysis and generation'),
    ('Django', 'Python framework for developing web applications'),
    ('Flask', 'Micro web framework for developing web applications with Python'),
    ('Data Science', 'Extracting knowledge and insights from data');

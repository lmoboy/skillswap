-- Course tables for SkillSwap

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

-- Sample data for courses
INSERT INTO courses (title, description, instructor_id, skill_id, difficulty_level, duration_hours, max_students, current_students, price, thumbnail_url, status) VALUES
('Introduction to Web Development', 'Learn the basics of HTML, CSS, and JavaScript to build your first website.', 1, 1, 'Beginner', 20, 30, 15, 49.99, 'https://images.unsplash.com/photo-1498050108023-c5249f4df085', 'Published'),
('Advanced React Patterns', 'Master advanced React patterns including hooks, context, and performance optimization.', 2, 2, 'Advanced', 15, 20, 8, 79.99, 'https://images.unsplash.com/photo-1633356122544-f134324a6cee', 'Published'),
('Python for Data Science', 'Learn Python programming with a focus on data analysis and visualization.', 3, 3, 'Intermediate', 25, 25, 12, 59.99, 'https://images.unsplash.com/photo-1526379095098-d400fd0bf935', 'Published'),
('UI/UX Design Fundamentals', 'Discover the principles of user interface and user experience design.', 8, 4, 'Beginner', 18, 40, 22, 39.99, 'https://images.unsplash.com/photo-1561070791-2526d30994b5', 'Published'),
('Mobile App Development with React Native', 'Build cross-platform mobile apps using React Native.', 6, 5, 'Intermediate', 30, 15, 7, 89.99, 'https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c', 'Published'),
('Machine Learning Basics', 'Introduction to machine learning algorithms and their applications.', 11, 6, 'Intermediate', 35, 20, 10, 99.99, 'https://images.unsplash.com/photo-1555949963-aa79dcee981c', 'Published'),
('DevOps and CI/CD Pipeline', 'Learn how to set up continuous integration and deployment pipelines.', 4, 7, 'Advanced', 22, 18, 9, 69.99, 'https://images.unsplash.com/photo-1618401471353-b98afee0b2eb', 'Published'),
('SEO Mastery Course', 'Master search engine optimization techniques to rank higher on Google.', 7, 8, 'Intermediate', 12, 50, 28, 44.99, 'https://images.unsplash.com/photo-1562577309-4932fdd64cd1', 'Published'),
('Graphic Design with Adobe Suite', 'Learn professional graphic design using Photoshop, Illustrator, and InDesign.', 44, 9, 'Beginner', 28, 30, 16, 54.99, 'https://images.unsplash.com/photo-1626785774573-4b799315345d', 'Published'),
('Full-Stack JavaScript Development', 'Become a full-stack developer using Node.js, Express, and MongoDB.', 9, 10, 'Advanced', 40, 25, 14, 109.99, 'https://images.unsplash.com/photo-1587620962725-abab7fe55159', 'Published');

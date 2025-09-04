CREATE DATABASE IF NOT EXISTS skillswap;
USE skillswap;


CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS skills (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE IF NOT EXISTS user_skills (
    user_id INT,
    skill_id INT,
    proficiency_level ENUM('Beginner', 'Intermediate', 'Advanced', 'Terry A Davis'),
    PRIMARY KEY (user_id, skill_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chats (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user1_id INT,
    user2_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    chat_id INT,
    sender_id INT,
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY chat_id REFERENCES chats(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
);



CREATE TABLE IF NOT EXISTS sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);



INSERT INTO users (username, email, password_hash) VALUES
('testuser', 'test@email.com', 'hashedpassword'),
('alice', 'alice@email.com', 'hashedpassword2'),
('bob', 'bob@email.com', 'hashedpassword3');

INSERT INTO skills (name, description) VALUES
('Python', 'Programming language for general-purpose development'),
('JavaScript', 'Language for web development'),
('SQL', 'Database query language'),
('Graphic Design', 'Creating visual content'),
('Public Speaking', 'Skill for effective oral communication');

INSERT INTO user_skills (user_id, skill_id, proficiency_level) VALUES
(1, 1, 'Advanced'),
(1, 2, 'Intermediate'),
(1, 3, 'Beginner');


INSERT INTO user_skills (user_id, skill_id, proficiency_level) VALUES
(2, 4, 'Advanced'),
(2, 2, 'Beginner'),
(3, 5, 'Intermediate');

INSERT INTO chats (user1_id, user2_id) VALUES
(1, 2),
(2, 3);

INSERT INTO messages (chat_id, sender_id, content) VALUES
(1,1, 'Hello Alice!'),
(1,2,'Hi Testuser, how can I help you?'),
(2,2,'Hi Bob!'),
(2,3,'Hello Alice!');
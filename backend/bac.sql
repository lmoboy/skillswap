DROP DATABASE IF EXISTS skillswap;
CREATE DATABASE IF NOT EXISTS skillswap;
USE skillswap;


CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
	 user_id VARCHAR(255) UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    profile_picture VARCHAR(255),
    aboutme TEXT,
    contacts TEXT,
    skills TEXT,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS skills (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE IF NOT EXISTS user_skills (
    user_id INT NOT null,
    skill_id INT NOT NULL,
    teaching_skill ENUM("Show how it's done", "Explain in slight details", "Give homework report on subject", "Professor") DEFAULT "Show how it's done",
    verified BOOL DEFAULT FALSE,
    PRIMARY KEY (user_id, skill_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chats (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user1_id INT NOT null,
    user2_id INT NOT null,
    initiated_by INT NOT null,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (initiated_by) REFERENCES users(id) ON DELETE CASCADE,
	 FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    chat_id INT NOT null,
    sender_id INT NOT null,
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT null,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO users (user_id, username, email, password_hash, profile_picture, aboutme, contacts) VALUES
("abcd1234efgh5678abcd1234efgh5678",'testuser', 'test@email.com', 'hashedpassword', 'testuser.jpg', 'Experienced developer specializing in web technologies.', '{"twitter":"@testuser", "linkedin":"linkedin.com/in/testuser"}'),
("bbcd1234efgh5678abcd1234efgh5678",'alice', 'alice@email.com', 'hashedpassword2', 'alice.png', 'Graphic designer with a passion for visual storytelling.', '{"portfolio":"alice-designs.com"}'),
("cbcd1234efgh5678abcd1234efgh5678",'bob', 'bob@email.com', 'hashedpassword3', 'bob.png', 'SQL and database expert.', '{"github":"github.com/bob-sql"}'),
("dbcd1234efgh5678abcd1234efgh5678",'charlie', 'charlie@email.com', 'hashedpassword4', 'charlie.png', 'Loves teaching public speaking to beginners.', '{"website":"charlie-speaking.com"}'),
("ebcd1234efgh5678abcd1234efgh5678",'diana', 'diana@email.com', 'hashedpassword5', 'diana.png', 'Front-end development enthusiast.', '{"twitter":"@diana_dev"}');

INSERT INTO skills (name, description) VALUES
('Python', 'Programming language for general-purpose development'),
('JavaScript', 'Language for web development'),
('SQL', 'Database query language'),
('Graphic Design', 'Creating visual content'),
('Public Speaking', 'Skill for effective oral communication'),
('React', 'JavaScript library for building user interfaces'),
('Data Science', 'Extracting knowledge and insights from data');

INSERT INTO user_skills (user_id, skill_id, teaching_skill, verified) VALUES
(1, 1, "Show how it's done", TRUE),
(1, 2, "Explain in slight details", FALSE),
(1, 3, "Give homework report on subject", TRUE),
(2, 4, "Professor", TRUE),
(3, 3, "Show how it's done", TRUE),
(4, 5, "Professor", FALSE),
(5, 2, "Show how it's done", TRUE),
(5, 6, "Explain in slight details", FALSE),
(5, 7, "Give homework report on subject", TRUE);

INSERT INTO chats (user1_id, user2_id, initiated_by) VALUES
(1, 2, 1),
(2, 3, 3),
(4, 5, 4),
(1, 4, 1);

INSERT INTO messages (chat_id, sender_id, content) VALUES
(1, 1, 'Hello Alice!'),
(1, 2, 'Hi Testuser, how can I help you?'),
(2, 2, 'Hi Bob!'),
(2, 3, 'Hello Alice!'),
(3, 4, 'Hey Diana, are you free to chat about front-end dev?'),
(3, 5, 'Sure Charlie, what do you need help with?'),
(4, 1, 'I need some help with public speaking.'),
(4, 4, 'No problem, what topic are you thinking of?');

SELECT
  u.*,
  s.name AS skill_name,
  s.description AS skill_description
FROM
  users AS u
JOIN
  user_skills AS us ON u.id = us.user_id
JOIN
  skills AS s ON us.skill_id = s.id
WHERE
  u.username LIKE '%test%'
  OR u.email LIKE '%test%'
  OR s.name LIKE '%web%' 
  OR s.description LIKE '%data%';
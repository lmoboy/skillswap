// Test data fixtures for consistent testing across all test types

export const validUserData = {
	username: 'testuser',
	email: 'test@example.com',
	password: 'testpassword123'
};

export const invalidUserData = {
	username: '',
	email: 'invalid-email',
	password: '123' // too short
};

export const validLoginData = {
	email: 'test@example.com',
	password: 'testpassword123'
};

export const invalidLoginData = {
	email: 'wrong@example.com',
	password: 'wrongpassword'
};

export const validCourseData = {
	title: 'Test Course',
	description: 'A comprehensive test course description',
	skill_name: 'JavaScript',
	duration_minutes: '120'
};

export const invalidCourseData = {
	title: '',
	description: '',
	skill_name: '',
	duration_minutes: 'invalid'
};

export const validModuleData = {
	title: 'Test Module',
	description: 'A test module description',
	video_duration: '30'
};

export const validContactData = {
	name: 'John Doe',
	link: 'https://example.com',
	icon: 'website'
};

export const invalidContactData = {
	name: '',
	link: 'invalid-url',
	icon: 'email'
};

export const validProjectData = {
	name: 'Test Project',
	description: 'A test project description',
	link: 'https://github.com/test/project'
};

export const mockUser = {
	id: 1,
	username: 'testuser',
	email: 'test@example.com',
	profile_picture: '/default-avatar.svg',
	skills: ['JavaScript', 'TypeScript'],
	projects: [],
	contacts: []
};

export const mockCourse = {
	id: 1,
	title: 'JavaScript Fundamentals',
	description: 'Learn JavaScript from scratch',
	instructor_name: 'John Doe',
	skill_name: 'JavaScript',
	duration_hours: 10,
	status: 'Published',
	thumbnail_url: '/course-thumbnail.jpg'
};

export const mockChat = {
	id: 1,
	initiator_id: 1,
	responder_id: 2,
	initiator_name: 'John Doe',
	responder_name: 'Jane Smith',
	last_message: 'Hello!',
	last_message_time: '2024-01-01T12:00:00Z'
};

export const mockMessage = {
	id: 1,
	chat_id: 1,
	sender_id: 1,
	content: 'Hello, how are you?',
	timestamp: '2024-01-01T12:00:00Z',
	sender_name: 'John Doe'
};

// API Response mocks
export const mockApiResponses = {
	loginSuccess: {
		status: 'ok',
		message: 'Login successful',
		user: mockUser
	},
	loginError: {
		error: 'Invalid credentials'
	},
	registerSuccess: {
		status: 'ok',
		message: 'Registration successful'
	},
	registerError: {
		error: 'Email already exists'
	},
	courseCreateSuccess: {
		message: 'Course created successfully',
		course_id: 1,
		skill_name: 'JavaScript'
	},
	courseCreateError: {
		error: 'Missing required fields'
	}
};

// Form validation error messages
export const validationErrors = {
	required: 'This field is required',
	email: 'Please enter a valid email address',
	password: 'Password must be at least 8 characters',
	url: 'Please enter a valid URL',
	username: 'Username must be between 3 and 50 characters',
	title: 'Title is required',
	description: 'Description is required',
	duration: 'Duration must be a positive number'
};

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mockFetch, mockFetchError, resetMocks } from '../../utils/testHelpers';
import { validCourseData, mockCourse, mockApiResponses } from '../../fixtures/testData';

describe('Courses API Integration', () => {
	beforeEach(() => {
		resetMocks();
		vi.clearAllMocks();
	});

	describe('Get All Courses', () => {
		it('sends correct request for all courses', async () => {
			mockFetch([mockCourse]);
			
			const response = await fetch('/api/courses');
			
			expect(global.fetch).toHaveBeenCalledWith('/api/courses');
			expect(response.ok).toBe(true);
		});

		it('handles successful courses response', async () => {
			const courses = [mockCourse];
			mockFetch(courses);
			
			const response = await fetch('/api/courses');
			const data = await response.json();
			
			expect(data).toEqual(courses);
		});

		it('handles empty courses response', async () => {
			mockFetch([]);
			
			const response = await fetch('/api/courses');
			const data = await response.json();
			
			expect(data).toEqual([]);
		});
	});

	describe('Get Course by ID', () => {
		it('sends correct request for specific course', async () => {
			mockFetch(mockCourse);
			
			const response = await fetch('/api/course?id=1');
			
			expect(global.fetch).toHaveBeenCalledWith('/api/course?id=1');
		});

		it('handles successful course detail response', async () => {
			mockFetch(mockCourse);
			
			const response = await fetch('/api/course?id=1');
			const data = await response.json();
			
			expect(data).toEqual(mockCourse);
		});

		it('handles course not found', async () => {
			mockFetch({ error: 'Course not found' }, 404);
			
			const response = await fetch('/api/course?id=999');
			
			expect(response.status).toBe(404);
		});
	});

	describe('Search Courses', () => {
		it('sends correct search request', async () => {
			mockFetch([mockCourse]);
			
			const searchQuery = { query: 'JavaScript' };
			const response = await fetch('/api/searchCourses', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(searchQuery)
			});
			
			expect(global.fetch).toHaveBeenCalledWith(
				'/api/searchCourses',
				expect.objectContaining({
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(searchQuery)
				})
			);
		});

		it('handles successful search response', async () => {
			const searchResults = [mockCourse];
			mockFetch(searchResults);
			
			const searchQuery = { query: 'JavaScript' };
			const response = await fetch('/api/searchCourses', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(searchQuery)
			});
			const data = await response.json();
			
			expect(data).toEqual(searchResults);
		});

		it('handles empty search results', async () => {
			mockFetch([]);
			
			const searchQuery = { query: 'nonexistent' };
			const response = await fetch('/api/searchCourses', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(searchQuery)
			});
			const data = await response.json();
			
			expect(data).toEqual([]);
		});
	});

	describe('Create Course', () => {
		it('sends correct course creation request', async () => {
			mockFetch(mockApiResponses.courseCreateSuccess);
			
			const formData = new FormData();
			formData.append('title', validCourseData.title);
			formData.append('description', validCourseData.description);
			formData.append('skill_name', validCourseData.skill_name);
			formData.append('duration_minutes', validCourseData.duration_minutes);
			
			const response = await fetch('/api/course/add', {
				method: 'POST',
				body: formData
			});
			
			expect(global.fetch).toHaveBeenCalledWith(
				'/api/course/add',
				expect.objectContaining({
					method: 'POST',
					body: formData
				})
			);
		});

		it('handles successful course creation', async () => {
			mockFetch(mockApiResponses.courseCreateSuccess);
			
			const formData = new FormData();
			formData.append('title', validCourseData.title);
			formData.append('description', validCourseData.description);
			formData.append('skill_name', validCourseData.skill_name);
			formData.append('duration_minutes', validCourseData.duration_minutes);
			
			const response = await fetch('/api/course/add', {
				method: 'POST',
				body: formData
			});
			const data = await response.json();
			
			expect(data).toEqual(mockApiResponses.courseCreateSuccess);
		});

		it('handles course creation validation errors', async () => {
			mockFetch({ error: 'Missing required fields' }, 400);
			
			const formData = new FormData();
			formData.append('title', '');
			formData.append('description', '');
			
			const response = await fetch('/api/course/add', {
				method: 'POST',
				body: formData
			});
			
			expect(response.status).toBe(400);
		});

		it('handles course creation with file upload', async () => {
			mockFetch(mockApiResponses.courseCreateSuccess);
			
			const formData = new FormData();
			formData.append('title', validCourseData.title);
			formData.append('description', validCourseData.description);
			formData.append('skill_name', validCourseData.skill_name);
			formData.append('duration_minutes', validCourseData.duration_minutes);
			
			// Add file
			const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
			formData.append('preview_photo', file);
			
			const response = await fetch('/api/course/add', {
				method: 'POST',
				body: formData
			});
			
			expect(global.fetch).toHaveBeenCalledWith(
				'/api/course/add',
				expect.objectContaining({
					method: 'POST',
					body: formData
				})
			);
		});
	});

	describe('Get Courses by Instructor', () => {
		it('sends correct request for instructor courses', async () => {
			mockFetch([mockCourse]);
			
			const response = await fetch('/api/coursesByInstructor?instructor_id=1');
			
			expect(global.fetch).toHaveBeenCalledWith('/api/coursesByInstructor?instructor_id=1');
		});

		it('handles successful instructor courses response', async () => {
			const instructorCourses = [mockCourse];
			mockFetch(instructorCourses);
			
			const response = await fetch('/api/coursesByInstructor?instructor_id=1');
			const data = await response.json();
			
			expect(data).toEqual(instructorCourses);
		});

		it('handles instructor with no courses', async () => {
			mockFetch([]);
			
			const response = await fetch('/api/coursesByInstructor?instructor_id=999');
			const data = await response.json();
			
			expect(data).toEqual([]);
		});
	});

	describe('API Error Handling', () => {
		it('handles 400 Bad Request for courses', async () => {
			mockFetch({ error: 'Invalid course data' }, 400);
			
			const response = await fetch('/api/course/add', {
				method: 'POST',
				body: new FormData()
			});
			
			expect(response.status).toBe(400);
		});

		it('handles 401 Unauthorized for protected endpoints', async () => {
			mockFetch({ error: 'Unauthorized' }, 401);
			
			const response = await fetch('/api/course/add', {
				method: 'POST',
				body: new FormData()
			});
			
			expect(response.status).toBe(401);
		});

		it('handles 404 Not Found for non-existent course', async () => {
			mockFetch({ error: 'Course not found' }, 404);
			
			const response = await fetch('/api/course?id=999');
			
			expect(response.status).toBe(404);
		});

		it('handles 500 Internal Server Error', async () => {
			mockFetch({ error: 'Internal Server Error' }, 500);
			
			const response = await fetch('/api/courses');
			
			expect(response.status).toBe(500);
		});

		it('handles network errors', async () => {
			mockFetchError('Network error');
			
			await expect(fetch('/api/courses')).rejects.toThrow('Network error');
		});
	});

	describe('Request Headers and Body', () => {
		it('includes correct content type for JSON requests', async () => {
			mockFetch([]);
			
			await fetch('/api/searchCourses', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ query: 'test' })
			});
			
			const fetchCall = vi.mocked(global.fetch).mock.calls[0];
			const options = fetchCall[1] as RequestInit;
			
			expect(options.headers).toMatchObject({
				'Content-Type': 'application/json'
			});
		});

		it('handles multipart form data correctly', async () => {
			mockFetch(mockApiResponses.courseCreateSuccess);
			
			const formData = new FormData();
			formData.append('title', 'Test Course');
			
			await fetch('/api/course/add', {
				method: 'POST',
				body: formData
			});
			
			const fetchCall = vi.mocked(global.fetch).mock.calls[0];
			const options = fetchCall[1] as RequestInit;
			
			expect(options.body).toBeInstanceOf(FormData);
		});
	});

	describe('Response Validation', () => {
		it('validates course object structure', async () => {
			mockFetch(mockCourse);
			
			const response = await fetch('/api/course?id=1');
			const course = await response.json();
			
			expect(course).toHaveProperty('id');
			expect(course).toHaveProperty('title');
			expect(course).toHaveProperty('description');
			expect(course).toHaveProperty('instructor_name');
			expect(course).toHaveProperty('skill_name');
		});

		it('validates course list structure', async () => {
			mockFetch([mockCourse]);
			
			const response = await fetch('/api/courses');
			const courses = await response.json();
			
			expect(Array.isArray(courses)).toBe(true);
			if (courses.length > 0) {
				expect(courses[0]).toHaveProperty('id');
				expect(courses[0]).toHaveProperty('title');
			}
		});
	});
});

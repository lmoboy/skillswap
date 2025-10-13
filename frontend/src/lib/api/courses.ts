/**
 * Courses API service
 */

const API_BASE = '/api';

const defaultOptions: RequestInit = {
    credentials: 'include',
    headers: {
        'Content-Type': 'application/json',
    },
};

/**
 * Fetch all courses
 */
export async function fetchCourses(): Promise<any[]> {
    try {
        const response = await fetch(`${API_BASE}/courses`, {
            ...defaultOptions,
            method: 'GET',
        });

        if (!response.ok) {
            throw new Error('Failed to fetch courses');
        }

        const data = await response.json();
        return data || [];
    } catch (error) {
        console.error('Error fetching courses:', error);
        throw error;
    }
}

/**
 * Fetch a single course by ID
 */
export async function fetchCourseById(courseId: string | number): Promise<any> {
    try {
        const response = await fetch(`${API_BASE}/courses/${courseId}`, {
            ...defaultOptions,
            method: 'GET',
        });

        if (!response.ok) {
            throw new Error('Failed to fetch course');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching course:', error);
        throw error;
    }
}

/**
 * Search courses by query
 */
export async function searchCourses(query: string): Promise<any[]> {
    try {
        const response = await fetch(`${API_BASE}/searchCourses`, {
            ...defaultOptions,
            method: 'POST',
            body: JSON.stringify({ query }),
        });

        if (!response.ok) {
            throw new Error('Failed to search courses');
        }

        const data = await response.json();
        return data || [];
    } catch (error) {
        console.error('Error searching courses:', error);
        throw error;
    }
}

/**
 * Create a new course
 */
export async function createCourse(formData: FormData): Promise<any> {
    try {
        const response = await fetch(`${API_BASE}/course/add`, {
            method: 'POST',
            body: formData,
            credentials: 'include',
            // Note: Don't set Content-Type header for FormData, browser will set it automatically
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ message: 'Server error' }));
            throw new Error(errorData.message || 'Failed to create course');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error creating course:', error);
        throw error;
    }
}

/**
 * Update a course
 */
export async function updateCourse(courseId: string | number, courseData: any): Promise<any> {
    try {
        const response = await fetch(`${API_BASE}/courses/${courseId}`, {
            ...defaultOptions,
            method: 'PUT',
            body: JSON.stringify(courseData),
        });

        if (!response.ok) {
            throw new Error('Failed to update course');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error updating course:', error);
        throw error;
    }
}

/**
 * Delete a course
 */
export async function deleteCourse(courseId: string | number): Promise<void> {
    try {
        const response = await fetch(`${API_BASE}/courses/${courseId}`, {
            ...defaultOptions,
            method: 'DELETE',
        });

        if (!response.ok) {
            throw new Error('Failed to delete course');
        }
    } catch (error) {
        console.error('Error deleting course:', error);
        throw error;
    }
}

/**
 * Enroll in a course
 */
export async function enrollInCourse(courseId: string | number): Promise<any> {
    try {
        const response = await fetch(`${API_BASE}/courses/${courseId}/enroll`, {
            ...defaultOptions,
            method: 'POST',
        });

        if (!response.ok) {
            throw new Error('Failed to enroll in course');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error enrolling in course:', error);
        throw error;
    }
}

/**
 * Get available skills
 */
export async function fetchSkills(): Promise<any[]> {
    try {
        const response = await fetch(`${API_BASE}/getSkills`, {
            ...defaultOptions,
            method: 'GET',
        });

        if (!response.ok) {
            throw new Error('Failed to fetch skills');
        }

        const data = await response.json();
        return data || [];
    } catch (error) {
        console.error('Error fetching skills:', error);
        throw error;
    }
}

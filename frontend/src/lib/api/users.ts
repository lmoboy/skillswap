/**
 * Users API service
 */

const API_BASE = '/api';

const defaultOptions: RequestInit = {
    credentials: 'include',
    headers: {
        'Content-Type': 'application/json',
    },
};

/**
 * Fetch user profile by ID
 */
export async function fetchUserProfile(userId: string | number): Promise<any> {
    try {
        const response = await fetch(`${API_BASE}/profile/${userId}`, {
            ...defaultOptions,
            method: 'GET',
        });

        if (!response.ok) {
            throw new Error('Failed to fetch user profile');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching user profile:', error);
        throw error;
    }
}

/**
 * Update user profile
 */
export async function updateUserProfile(userId: string | number, profileData: any): Promise<any> {
    try {
        const response = await fetch(`${API_BASE}/profile/${userId}`, {
            ...defaultOptions,
            method: 'PUT',
            body: JSON.stringify(profileData),
        });

        if (!response.ok) {
            throw new Error('Failed to update profile');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error updating profile:', error);
        throw error;
    }
}

/**
 * Upload user profile picture
 */
export async function uploadProfilePicture(userId: string | number, file: File): Promise<any> {
    try {
        const formData = new FormData();
        formData.append('picture', file);

        const response = await fetch(`${API_BASE}/profile/${userId}/picture`, {
            method: 'POST',
            body: formData,
            credentials: 'include',
        });

        if (!response.ok) {
            throw new Error('Failed to upload profile picture');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error uploading profile picture:', error);
        throw error;
    }
}

/**
 * Search users by query
 */
export async function searchUsers(query: string): Promise<any[]> {
    try {
        const response = await fetch(`${API_BASE}/search`, {
            ...defaultOptions,
            method: 'POST',
            body: JSON.stringify({ query }),
        });

        if (!response.ok) {
            throw new Error('Failed to search users');
        }

        const data = await response.json();
        return data || [];
    } catch (error) {
        console.error('Error searching users:', error);
        throw error;
    }
}

/**
 * Full search (users and courses)
 */
export async function fullSearch(query: string): Promise<any> {
    try {
        const response = await fetch(`${API_BASE}/fullSearch`, {
            ...defaultOptions,
            method: 'POST',
            body: JSON.stringify({ query }),
        });

        if (!response.ok) {
            throw new Error('Failed to perform full search');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error performing full search:', error);
        throw error;
    }
}

/**
 * Add skill to user profile
 */
export async function addSkill(userId: string | number, skillName: string): Promise<any> {
    try {
        const response = await fetch(`${API_BASE}/profile/${userId}/skills`, {
            ...defaultOptions,
            method: 'POST',
            body: JSON.stringify({ skill_name: skillName }),
        });

        if (!response.ok) {
            throw new Error('Failed to add skill');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error adding skill:', error);
        throw error;
    }
}

/**
 * Remove skill from user profile
 */
export async function removeSkill(userId: string | number, skillName: string): Promise<void> {
    try {
        const response = await fetch(`${API_BASE}/profile/${userId}/skills/${skillName}`, {
            ...defaultOptions,
            method: 'DELETE',
        });

        if (!response.ok) {
            throw new Error('Failed to remove skill');
        }
    } catch (error) {
        console.error('Error removing skill:', error);
        throw error;
    }
}

/**
 * Add project to user profile
 */
export async function addProject(userId: string | number, project: any): Promise<any> {
    try {
        const response = await fetch(`${API_BASE}/profile/${userId}/projects`, {
            ...defaultOptions,
            method: 'POST',
            body: JSON.stringify(project),
        });

        if (!response.ok) {
            throw new Error('Failed to add project');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error adding project:', error);
        throw error;
    }
}

/**
 * Remove project from user profile
 */
export async function removeProject(userId: string | number, projectId: string | number): Promise<void> {
    try {
        const response = await fetch(`${API_BASE}/profile/${userId}/projects/${projectId}`, {
            ...defaultOptions,
            method: 'DELETE',
        });

        if (!response.ok) {
            throw new Error('Failed to remove project');
        }
    } catch (error) {
        console.error('Error removing project:', error);
        throw error;
    }
}

/**
 * Add contact to user profile
 */
export async function addContact(userId: string | number, contact: any): Promise<any> {
    try {
        const response = await fetch(`${API_BASE}/profile/${userId}/contacts`, {
            ...defaultOptions,
            method: 'POST',
            body: JSON.stringify(contact),
        });

        if (!response.ok) {
            throw new Error('Failed to add contact');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error adding contact:', error);
        throw error;
    }
}

/**
 * Remove contact from user profile
 */
export async function removeContact(userId: string | number, contactId: string | number): Promise<void> {
    try {
        const response = await fetch(`${API_BASE}/profile/${userId}/contacts/${contactId}`, {
            ...defaultOptions,
            method: 'DELETE',
        });

        if (!response.ok) {
            throw new Error('Failed to remove contact');
        }
    } catch (error) {
        console.error('Error removing contact:', error);
        throw error;
    }
}

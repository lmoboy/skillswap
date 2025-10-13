/**
 * Validation utility functions
 */

/**
 * Validates if a string is a valid email address
 */
export function isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

/**
 * Validates if a string is a valid URL
 */
export function isValidUrl(url: string): boolean {
    try {
        new URL(url);
        return true;
    } catch {
        return false;
    }
}

/**
 * Validates password strength
 * Returns error message if invalid, null if valid
 */
export function validatePassword(password: string): string | null {
    if (password.length < 8) {
        return 'Password must be at least 8 characters long';
    }
    if (password.length > 50) {
        return 'Password cannot be longer than 50 characters';
    }
    return null;
}

/**
 * Validates username
 * Returns error message if invalid, null if valid
 */
export function validateUsername(username: string): string | null {
    if (!username || username.trim().length === 0) {
        return 'Username is required';
    }
    if (username.length > 50) {
        return 'Username is too long (max 50 characters)';
    }
    return null;
}

/**
 * Validates email
 * Returns error message if invalid, null if valid
 */
export function validateEmail(email: string): string | null {
    if (!email || email.trim().length === 0) {
        return 'Email is required';
    }
    if (email.length > 100) {
        return 'Email is too long (max 100 characters)';
    }
    if (!isValidEmail(email)) {
        return 'Please enter a valid email address';
    }
    return null;
}

/**
 * Validates a required field
 */
export function validateRequired(value: string, fieldName: string): string | null {
    if (!value || value.trim().length === 0) {
        return `${fieldName} is required`;
    }
    return null;
}

/**
 * Validates maximum length
 */
export function validateMaxLength(value: string, maxLength: number, fieldName: string): string | null {
    if (value.length > maxLength) {
        return `${fieldName} cannot exceed ${maxLength} characters`;
    }
    return null;
}

/**
 * Validates minimum length
 */
export function validateMinLength(value: string, minLength: number, fieldName: string): string | null {
    if (value.length < minLength) {
        return `${fieldName} must be at least ${minLength} characters`;
    }
    return null;
}

/**
 * Validates a number is positive
 */
export function validatePositiveNumber(value: number, fieldName: string): string | null {
    if (value <= 0) {
        return `${fieldName} must be a positive number`;
    }
    return null;
}

/**
 * Validates file type
 */
export function validateFileType(file: File, allowedTypes: string[]): string | null {
    if (!allowedTypes.includes(file.type)) {
        return `File type not allowed. Allowed types: ${allowedTypes.join(', ')}`;
    }
    return null;
}

/**
 * Validates file size (in bytes)
 */
export function validateFileSize(file: File, maxSize: number): string | null {
    if (file.size > maxSize) {
        const maxSizeMB = (maxSize / (1024 * 1024)).toFixed(2);
        return `File size exceeds ${maxSizeMB}MB`;
    }
    return null;
}

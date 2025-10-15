import { describe, it, expect } from 'vitest';
import { 
	formatDate, 
	formatDuration, 
	formatFileSize, 
	truncateText,
	formatRelativeTime,
	formatTime,
	formatNumber,
	toTitleCase,
	formatRating,
	pluralize,
	getInitials,
	formatPrice,
	formatVideoTime
} from '$lib/utils/formatting';

describe('Formatting Utilities', () => {
	describe('formatDate', () => {
		it('formats date string correctly', () => {
			const dateString = '2024-01-15T10:30:00Z';
			const formatted = formatDate(dateString);
			expect(formatted).toMatch(/January 15, 2024/);
		});

		it('handles invalid date gracefully', () => {
			expect(formatDate('invalid-date')).toBe('Invalid Date');
		});
	});

	describe('formatDuration', () => {
		it('formats minutes correctly', () => {
			expect(formatDuration(30)).toBe('30m');
			expect(formatDuration(60)).toBe('1h');
			expect(formatDuration(90)).toBe('1h 30m');
		});

		it('handles zero duration', () => {
			expect(formatDuration(0)).toBe('0m');
		});
	});

	describe('formatFileSize', () => {
		it('formats bytes correctly', () => {
			expect(formatFileSize(0)).toBe('0 Bytes');
			expect(formatFileSize(500)).toBe('500 Bytes');
			expect(formatFileSize(1024)).toBe('1 KB');
		});

		it('formats larger sizes correctly', () => {
			expect(formatFileSize(1048576)).toBe('1 MB');
			expect(formatFileSize(1073741824)).toBe('1 GB');
		});
	});

	describe('truncateText', () => {
		it('truncates text to specified length', () => {
			const text = 'This is a very long text that should be truncated';
			expect(truncateText(text, 20)).toBe('This is a very lo...');
		});

		it('returns original text if shorter than limit', () => {
			const text = 'Short text';
			expect(truncateText(text, 20)).toBe('Short text');
		});

		it('handles empty string', () => {
			expect(truncateText('', 10)).toBe('');
		});
	});

	describe('formatRelativeTime', () => {
		it('formats recent times correctly', () => {
			const now = new Date().toISOString();
			expect(formatRelativeTime(now)).toBe('just now');
		});
	});

	describe('formatTime', () => {
		it('formats time correctly', () => {
			const timeString = '2024-01-15T10:30:00Z';
			const formatted = formatTime(timeString);
			expect(formatted).toMatch(/\d{1,2}:\d{2}/);
		});
	});

	describe('formatNumber', () => {
		it('formats numbers with commas', () => {
			expect(formatNumber(1000)).toBe('1,000');
			expect(formatNumber(1000000)).toBe('1,000,000');
		});
	});

	describe('toTitleCase', () => {
		it('converts string to title case', () => {
			expect(toTitleCase('hello world')).toBe('Hello World');
			expect(toTitleCase('test string')).toBe('Test String');
		});
	});

	describe('formatRating', () => {
		it('formats rating with one decimal place', () => {
			expect(formatRating(4.567)).toBe('4.6');
			expect(formatRating(5)).toBe('5.0');
		});
	});

	describe('pluralize', () => {
		it('handles singular and plural forms', () => {
			expect(pluralize(1, 'item')).toBe('item');
			expect(pluralize(2, 'item')).toBe('items');
		});
	});

	describe('getInitials', () => {
		it('generates initials from name', () => {
			expect(getInitials('John Doe')).toBe('JD');
			expect(getInitials('Jane')).toBe('JA');
		});
	});

	describe('formatPrice', () => {
		it('formats price with currency', () => {
			expect(formatPrice(19.99)).toBe('$19.99');
			expect(formatPrice(100, 'EUR')).toBe('€100.00');
		});
	});

	describe('formatVideoTime', () => {
		it('formats seconds to MM:SS format', () => {
			expect(formatVideoTime(65)).toBe('1:05');
			expect(formatVideoTime(125)).toBe('2:05');
		});
	});
});
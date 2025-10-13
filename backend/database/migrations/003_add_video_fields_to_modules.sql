-- Migration: Add video-related fields to course_modules table
-- Date: 2024
-- Description: Adds video_url, video_duration, and thumbnail_url fields to support video modules

ALTER TABLE course_modules
ADD COLUMN video_url VARCHAR(500) AFTER description,
ADD COLUMN video_duration INT UNSIGNED DEFAULT 0 AFTER video_url,
ADD COLUMN thumbnail_url VARCHAR(500) AFTER video_duration;

-- Update existing modules to have default values
UPDATE course_modules
SET video_url = '',
    video_duration = 0,
    thumbnail_url = ''
WHERE video_url IS NULL;

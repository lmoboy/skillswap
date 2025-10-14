-- Migration: Add video-related fields to course_modules table
-- Date: 2024
-- Description: Adds video_url, video_duration, and thumbnail_url fields to support video modules
-- Note: This migration is safe to run multiple times - it only adds columns if they don't exist

-- Check if columns exist before adding them
SET @dbname = DATABASE();
SET @tablename = 'course_modules';

-- Add video_url if it doesn't exist
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
  WHERE TABLE_SCHEMA = @dbname AND TABLE_NAME = @tablename AND COLUMN_NAME = 'video_url');

SET @query = IF(@col_exists = 0, 
  'ALTER TABLE course_modules ADD COLUMN video_url VARCHAR(500) AFTER description', 
  'SELECT "Column video_url already exists" AS msg');

PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add video_duration if it doesn't exist
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
  WHERE TABLE_SCHEMA = @dbname AND TABLE_NAME = @tablename AND COLUMN_NAME = 'video_duration');

SET @query = IF(@col_exists = 0, 
  'ALTER TABLE course_modules ADD COLUMN video_duration INT UNSIGNED DEFAULT 0 AFTER video_url', 
  'SELECT "Column video_duration already exists" AS msg');

PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add thumbnail_url if it doesn't exist
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
  WHERE TABLE_SCHEMA = @dbname AND TABLE_NAME = @tablename AND COLUMN_NAME = 'thumbnail_url');

SET @query = IF(@col_exists = 0, 
  'ALTER TABLE course_modules ADD COLUMN thumbnail_url VARCHAR(500) AFTER video_duration', 
  'SELECT "Column thumbnail_url already exists" AS msg');

PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;


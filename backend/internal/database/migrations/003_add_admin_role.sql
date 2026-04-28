-- Add admin role support to users table
-- Migration: 003_add_admin_role.sql

-- Check if column exists before adding it
SET @dbname = DATABASE();
SET @tablename = 'users';
SET @columnname = 'is_admin';

SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
  WHERE TABLE_SCHEMA = @dbname AND TABLE_NAME = @tablename AND COLUMN_NAME = @columnname);

SET @query = IF(@col_exists = 0, 
  'ALTER TABLE users ADD COLUMN is_admin TINYINT(1) NOT NULL DEFAULT 0 AFTER swaps', 
  'SELECT "Column is_admin already exists" AS msg');

PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add index for admin lookups if it doesn't exist
SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
  WHERE TABLE_SCHEMA = @dbname AND TABLE_NAME = @tablename AND INDEX_NAME = 'idx_users_is_admin');

SET @query = IF(@index_exists = 0,
  'CREATE INDEX idx_users_is_admin ON users(is_admin)',
  'SELECT "Index idx_users_is_admin already exists" AS msg');

PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Optional: Set a specific user as admin (uncomment and modify email as needed)
-- UPDATE users SET is_admin = 1 WHERE email = 'admin@example.com';

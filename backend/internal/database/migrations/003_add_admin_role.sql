-- Add admin role support to users table
-- Migration: 003_add_admin_role.sql

-- Add is_admin column to users table
ALTER TABLE users
ADD COLUMN is_admin TINYINT(1) NOT NULL DEFAULT 0 AFTER swaps;

-- Add index for admin lookups
CREATE INDEX idx_users_is_admin ON users(is_admin);

-- Optional: Set a specific user as admin (uncomment and modify email as needed)
-- UPDATE users SET is_admin = 1 WHERE email = 'admin@example.com';

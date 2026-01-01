-- Migration: Add email verification and trial fields
-- Version: 003
-- Description: Add email verification token, verified timestamp, and trial end date to users table

-- Add email verification fields
ALTER TABLE users 
ADD COLUMN email_verified_at TIMESTAMP,
ADD COLUMN verification_token VARCHAR(255),
ADD COLUMN verification_token_expires_at TIMESTAMP,
ADD COLUMN trial_ends_at TIMESTAMP;

-- Create index on verification token for faster lookups
CREATE INDEX idx_users_verification_token ON users(verification_token) WHERE verification_token IS NOT NULL;

-- Create index on trial_ends_at for cron job queries
CREATE INDEX idx_users_trial_ends_at ON users(trial_ends_at) WHERE trial_ends_at IS NOT NULL;

-- Update existing users to have email_verified = true (migration safety)
UPDATE users SET email_verified = true WHERE email_verified_at IS NULL;
UPDATE users SET email_verified_at = created_at WHERE email_verified = true AND email_verified_at IS NULL;

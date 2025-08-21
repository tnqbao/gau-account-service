-- Remove verification columns as they are now handled by user_verifications table
ALTER TABLE users DROP COLUMN IF EXISTS is_phone_verified;
ALTER TABLE users DROP COLUMN IF EXISTS is_email_verified;

-- Update constraints to match new entity structure
DROP INDEX IF EXISTS uq_users_email_permission;
DROP INDEX IF EXISTS uq_users_phone_permission;

-- Add new unique constraints for email, username, facebook_url, github_url as per entity
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX idx_users_facebook_url ON users(facebook_url) WHERE facebook_url IS NOT NULL;
CREATE UNIQUE INDEX idx_users_github_url ON users(github_url) WHERE github_url IS NOT NULL;

-- Keep the existing username_permission composite index
CREATE INDEX IF NOT EXISTS idx_users_username_permission ON users(username, permission);

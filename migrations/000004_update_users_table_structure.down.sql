-- Reverse the changes made in the up migration

-- Drop the new unique indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_facebook_url;
DROP INDEX IF EXISTS idx_users_github_url;

-- Add back the verification columns
ALTER TABLE users ADD COLUMN is_phone_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN is_email_verified BOOLEAN DEFAULT FALSE;

-- Restore the original constraints
ALTER TABLE users ADD CONSTRAINT uq_users_email_permission UNIQUE (email, permission);
ALTER TABLE users ADD CONSTRAINT uq_users_phone_permission UNIQUE (phone, permission);

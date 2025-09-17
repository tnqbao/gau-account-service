-- Final consolidated migration for gau-account-service database schema
-- This migration creates the complete database structure

-- Drop existing tables if they exist
DROP TABLE IF EXISTS user_mfas;
DROP TABLE IF EXISTS user_verifications;
DROP TABLE IF EXISTS users;

-- Create users table with final structure
CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    full_name VARCHAR(255),
    gender VARCHAR(10),
    date_of_birth DATE,
    facebook_url VARCHAR(255),
    github_url VARCHAR(255),
    permission VARCHAR(255),
    username VARCHAR(255),
    password VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(15),
    avatar_url VARCHAR(255) DEFAULT 'https://cdn.gauas.online/images/avatar/default_image.jpg',

    CONSTRAINT uq_users_username_permission UNIQUE (username, permission)
);

-- Create indexes for users table
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX idx_users_facebook_url ON users(facebook_url) WHERE facebook_url IS NOT NULL;
CREATE UNIQUE INDEX idx_users_github_url ON users(github_url) WHERE github_url IS NOT NULL;
CREATE INDEX idx_users_username_permission ON users(username, permission);

-- Create user_verifications table
CREATE TABLE user_verifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    method VARCHAR(20) NOT NULL,
    value VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_verifications_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Create indexes for user_verifications table
CREATE INDEX idx_user_verifications_user_id ON user_verifications(user_id);
CREATE INDEX idx_user_verifications_method ON user_verifications(method);

-- Create user_mfas table
CREATE TABLE user_mfas (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type VARCHAR(30) NOT NULL,
    secret VARCHAR(255),
    enabled BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_mfas_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Create indexes for user_mfas table
CREATE INDEX idx_user_mfas_user_id ON user_mfas(user_id);
CREATE INDEX idx_user_mfas_type ON user_mfas(type);

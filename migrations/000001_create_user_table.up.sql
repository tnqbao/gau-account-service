DROP TABLE IF EXISTS users;

CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    full_name VARCHAR(255),
    gender VARCHAR(10),
    date_of_birth DATE,
    facebook_url VARCHAR(255) UNIQUE,
    github_url VARCHAR(255) UNIQUE,
    permission VARCHAR(255),
    username VARCHAR(255),
    password VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(15),
    is_phone_verified BOOLEAN DEFAULT FALSE,
    is_email_verified BOOLEAN DEFAULT FALSE,

    CONSTRAINT uq_users_username_permission UNIQUE (username, permission),
    CONSTRAINT uq_users_email_permission UNIQUE (email, permission)
);

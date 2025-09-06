-- Rollback migration for final consolidated schema
-- This migration drops all tables created by the final schema

DROP TABLE IF EXISTS user_mfas;
DROP TABLE IF EXISTS user_verifications;
DROP TABLE IF EXISTS users;

-- Drop all triggers first
DROP TRIGGER IF EXISTS update_users_modified_at ON users;
DROP TRIGGER IF EXISTS update_system_package_filters_modified_at ON system_package_filters;
DROP TRIGGER IF EXISTS update_sandboxes_modified_at ON sandboxes;
DROP TRIGGER IF EXISTS update_executions_modified_at ON executions;
DROP TRIGGER IF EXISTS update_jobs_modified_at ON jobs;
DROP TRIGGER IF EXISTS update_workers_modified_at ON workers;
DROP TRIGGER IF EXISTS update_exec_request_modified_at ON exec_request;
DROP TRIGGER IF EXISTS update_language_versions_modified_at ON language_versions;
DROP TRIGGER IF EXISTS update_languages_modified_at ON languages;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_modified_at_column();

-- Recreate all tables without inheritance to restore original structure

-- Drop existing tables (in reverse dependency order)
DROP TABLE IF EXISTS system_package_filters CASCADE;
DROP TABLE IF EXISTS sandboxes CASCADE;
DROP TABLE IF EXISTS executions CASCADE;
DROP TABLE IF EXISTS jobs CASCADE;
DROP TABLE IF EXISTS workers CASCADE;
DROP TABLE IF EXISTS exec_request CASCADE;
DROP TABLE IF EXISTS language_versions CASCADE;
DROP TABLE IF EXISTS languages CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop the base audit mixin table
DROP TABLE IF EXISTS fullauditmixin CASCADE;

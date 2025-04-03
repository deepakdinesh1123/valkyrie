-- Drop executions table and related sequence
drop table if exists executions;
drop sequence if exists executions_id_seq;

-- Drop jobs table and related sequence
drop table if exists jobs;
drop sequence if exists jobs_id_seq;

-- Drop workers table and related sequence
drop table if exists workers;
drop sequence if exists workers_id_seq;

-- Drop job_types table
drop table if exists job_types;

-- Drop job_groups table
drop table if exists job_groups;

-- Drop exec_request table and related sequence
drop table if exists exec_request;
drop sequence if exists exec_request_id_seq;

-- Drop language_versions table and related index and sequence
drop index if exists unique_default_version_per_language;
drop table if exists language_versions;
drop sequence if exists language_versions_id_seq;

-- Drop languages table and related sequence
drop table if exists languages;
drop sequence if exists languages_id_seq;

-- Drop packages table and related sequence
drop table if exists packages;
drop sequence if exists packages_id_seq;

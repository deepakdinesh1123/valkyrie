drop table if exists job_groups;
drop table if exists job_types;
drop table if exists workers;
drop sequence if exists workers_id_seq;
drop table if exists jobs;
drop sequence if exists jobs_id_seq;
drop table if exists job_runs;
drop sequence if exists job_runs_id_seq;
drop table if exists execute;
drop sequence if exists execute_id_seq;

-- Drop the foreign key constraint and the table `language_versions` first
DROP TABLE IF EXISTS language_versions CASCADE;

-- Drop the sequence for `language_versions`
DROP SEQUENCE IF EXISTS language_versions_id_seq;

-- Drop the table `languages`
DROP TABLE IF EXISTS languages CASCADE;

-- Drop the sequence for `languages`
DROP SEQUENCE IF EXISTS languages_id_seq;
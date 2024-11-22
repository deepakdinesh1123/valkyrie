create sequence languages_id_seq as bigint;

CREATE TABLE languages (
    id bigint PRIMARY KEY DEFAULT nextval('languages_id_seq'),
    name TEXT NOT NULL UNIQUE,                  
    extension TEXT NOT NULL,                    
    monaco_language TEXT NOT NULL,
    template TEXT NOT NULL,
    default_code TEXT NOT NULL                    
);

create sequence language_versions_id_seq as bigint;

CREATE TABLE language_versions (
    id bigint PRIMARY KEY DEFAULT nextval('language_versions_id_seq'),
    language_id BIGINT NOT NULL REFERENCES languages (id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    nix_package_name TEXT NOT NULL,
    template TEXT,
    search_query TEXT NOT NULL, 
    default_version BOOLEAN NOT NULL DEFAULT false,                          
    UNIQUE (language_id, version)               
);

CREATE UNIQUE INDEX unique_default_version_per_language 
ON language_versions (language_id) 
WHERE default_version = true;

-- name: TruncateLanguages :exec
TRUNCATE TABLE languages CASCADE;

-- name: TruncateLanguageVersions :exec
TRUNCATE TABLE language_versions CASCADE;

-- name: CreateLanguage :one
INSERT INTO languages (name, extension, monaco_language, default_code) 
VALUES ($1, $2, $3, $4) 
RETURNING id;

-- name: GetLanguageByID :one
SELECT *
FROM languages 
WHERE id = $1;

-- name: GetAllLanguages :many
SELECT *
FROM languages;

-- name: UpdateLanguage :one
UPDATE languages 
SET name = $2, extension = $3, monaco_language = $4, default_code = $5
WHERE id = $1
returning id;

-- name: DeleteLanguage :one
DELETE FROM languages 
WHERE id = $1 
returning id;

-- name: CreateLanguageVersion :one
INSERT INTO language_versions (
    language_id, version, nix_package_name, template, search_query, default_version
) 
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING id;

-- name: GetAllLanguageVersions :many
SELECT id, language_id, version, nix_package_name, template, search_query, default_version
FROM language_versions;

-- name: GetLanguageVersionByID :one
SELECT id, language_id, version, nix_package_name, template, search_query, default_version 
FROM language_versions 
WHERE id = $1;

-- name: GetVersionsByLanguageID :many
SELECT id, language_id, version, nix_package_name, template, search_query, default_version 
FROM language_versions 
WHERE language_id = $1;

-- name: GetLanguageVersion :one
SELECT id, language_id, version, nix_package_name, template, search_query, default_version 
FROM language_versions 
WHERE language_id = $1 AND version = $2;

-- name: UpdateLanguageVersion :one
UPDATE language_versions 
SET language_id = $2, nix_package_name = $3, template = $4, search_query = $5, version = $6, default_version = $7
WHERE id = $1 
returning id;

-- name: DeleteLanguageVersion :one
DELETE FROM language_versions 
WHERE id = $1
returning language_id;

-- name: DeleteAllVersionsForLanguage :one
DELETE FROM language_versions 
WHERE language_id = $1
returning id;

-- name: GetLanguageByName :one
SELECT * from languages WHERE name = $1;

-- name: GetDefaultVersion :one
SELECT * FROM language_versions WHERE default_version = true AND language_id = $1;

-- name: InsertLanguages :copyfrom
INSERT INTO languages (name, extension, monaco_language, default_code, template) VALUES ($1, $2, $3, $4, $5);

-- name: InsertLanguageVersions :copyfrom
INSERT INTO language_versions (language_id, version, nix_package_name, template, search_query, default_version) VALUES ($1, $2, $3, $4, $5, $6);
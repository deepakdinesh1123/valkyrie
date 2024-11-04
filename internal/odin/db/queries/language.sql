CREATE TABLE languages (
    id bigint PRIMARY KEY DEFAULT nextval('languages_id_seq'),
    name TEXT NOT NULL UNIQUE,                  
    extension TEXT NOT NULL,
    monaco_language TEXT NOT NULL,                      
    default_code TEXT NOT NULL                           
);

CREATE TABLE language_versions (
    id bigint PRIMARY KEY DEFAULT nextval('language_versions_id_seq'),
    language_id BIGINT NOT NULL REFERENCES languages (id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    nix_package_name TEXT NOT NULL,             
    flake_template TEXT NOT NULL,
    script_template TEXT NOT NULL,                       
    search_query TEXT NOT NULL,                           
    UNIQUE (language_id, version)               
);

-- name: CreateLanguage :one
INSERT INTO languages (name, extension, monaco_language, default_code) 
VALUES ($1, $2, $3, $4) 
RETURNING id;

-- name: GetLanguageByID :one
SELECT id, name, extension, monaco_language, default_code 
FROM languages 
WHERE id = $1;

-- name: GetAllLanguages :many
SELECT id, name, extension, monaco_language, default_code 
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
    language_id, version, nix_package_name, flake_template, script_template, search_query
) 
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING id;

-- name: GetAllLanguageVersions :many
SELECT id, language_id, version, nix_package_name, flake_template, script_template, search_query 
FROM language_versions;

-- name: GetLanguageVersionByID :one
SELECT id, language_id, version, nix_package_name, flake_template, script_template, search_query  
FROM language_versions 
WHERE id = $1;

-- name: GetVersionsByLanguageID :many
SELECT id, language_id, version, nix_package_name, flake_template, script_template, search_query 
FROM language_versions 
WHERE language_id = $1;

-- name: GetLanguageVersion :one
SELECT id, language_id, version, nix_package_name, flake_template, script_template, search_query 
FROM language_versions 
WHERE language_id = $1 AND version = $2;

-- name: UpdateLanguageVersion :one
UPDATE language_versions 
SET language_id = $2, nix_package_name = $3, flake_template = $4, script_template = $5, search_query = $6, version = $7
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
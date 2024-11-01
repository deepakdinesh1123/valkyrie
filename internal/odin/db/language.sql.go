// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: language.sql

package db

import (
	"context"
)

const createLanguage = `-- name: CreateLanguage :one
INSERT INTO languages (name, extension, monaco_language) 
VALUES ($1, $2, $3) 
RETURNING id
`

type CreateLanguageParams struct {
	Name           string `db:"name" json:"name"`
	Extension      string `db:"extension" json:"extension"`
	MonacoLanguage string `db:"monaco_language" json:"monaco_language"`
}

func (q *Queries) CreateLanguage(ctx context.Context, arg CreateLanguageParams) (int64, error) {
	row := q.db.QueryRow(ctx, createLanguage, arg.Name, arg.Extension, arg.MonacoLanguage)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createLanguageVersion = `-- name: CreateLanguageVersion :one
INSERT INTO language_versions (
    language_id, version, nix_package_name, flake_template, script_template, default_code, search_query
) 
VALUES ($1, $2, $3, $4, $5, $6, $7) 
RETURNING id
`

type CreateLanguageVersionParams struct {
	LanguageID     int64  `db:"language_id" json:"language_id"`
	Version        string `db:"version" json:"version"`
	NixPackageName string `db:"nix_package_name" json:"nix_package_name"`
	FlakeTemplate  string `db:"flake_template" json:"flake_template"`
	ScriptTemplate string `db:"script_template" json:"script_template"`
	DefaultCode    string `db:"default_code" json:"default_code"`
	SearchQuery    string `db:"search_query" json:"search_query"`
}

func (q *Queries) CreateLanguageVersion(ctx context.Context, arg CreateLanguageVersionParams) (int64, error) {
	row := q.db.QueryRow(ctx, createLanguageVersion,
		arg.LanguageID,
		arg.Version,
		arg.NixPackageName,
		arg.FlakeTemplate,
		arg.ScriptTemplate,
		arg.DefaultCode,
		arg.SearchQuery,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteAllVersionsForLanguage = `-- name: DeleteAllVersionsForLanguage :one
DELETE FROM language_versions 
WHERE language_id = $1
returning id
`

func (q *Queries) DeleteAllVersionsForLanguage(ctx context.Context, languageID int64) (int64, error) {
	row := q.db.QueryRow(ctx, deleteAllVersionsForLanguage, languageID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteLanguage = `-- name: DeleteLanguage :one
DELETE FROM languages 
WHERE id = $1 
returning id
`

func (q *Queries) DeleteLanguage(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRow(ctx, deleteLanguage, id)
	err := row.Scan(&id)
	return id, err
}

const deleteLanguageVersion = `-- name: DeleteLanguageVersion :one
DELETE FROM language_versions 
WHERE id = $1
returning language_id
`

func (q *Queries) DeleteLanguageVersion(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRow(ctx, deleteLanguageVersion, id)
	var language_id int64
	err := row.Scan(&language_id)
	return language_id, err
}

const getAllLanguageVersions = `-- name: GetAllLanguageVersions :many
SELECT id, language_id, version, nix_package_name, flake_template, script_template, default_code, search_query 
FROM language_versions
`

func (q *Queries) GetAllLanguageVersions(ctx context.Context) ([]LanguageVersion, error) {
	rows, err := q.db.Query(ctx, getAllLanguageVersions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LanguageVersion
	for rows.Next() {
		var i LanguageVersion
		if err := rows.Scan(
			&i.ID,
			&i.LanguageID,
			&i.Version,
			&i.NixPackageName,
			&i.FlakeTemplate,
			&i.ScriptTemplate,
			&i.DefaultCode,
			&i.SearchQuery,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllLanguages = `-- name: GetAllLanguages :many
SELECT id, name, extension, monaco_language 
FROM languages
`

func (q *Queries) GetAllLanguages(ctx context.Context) ([]Language, error) {
	rows, err := q.db.Query(ctx, getAllLanguages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Language
	for rows.Next() {
		var i Language
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Extension,
			&i.MonacoLanguage,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLanguageByID = `-- name: GetLanguageByID :one
SELECT id, name, extension, monaco_language 
FROM languages 
WHERE id = $1
`

func (q *Queries) GetLanguageByID(ctx context.Context, id int64) (Language, error) {
	row := q.db.QueryRow(ctx, getLanguageByID, id)
	var i Language
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Extension,
		&i.MonacoLanguage,
	)
	return i, err
}

const getLanguageVersion = `-- name: GetLanguageVersion :one
SELECT id, language_id, version, nix_package_name, flake_template, script_template, default_code, search_query 
FROM language_versions 
WHERE language_id = $1 AND version = $2
`

type GetLanguageVersionParams struct {
	LanguageID int64  `db:"language_id" json:"language_id"`
	Version    string `db:"version" json:"version"`
}

func (q *Queries) GetLanguageVersion(ctx context.Context, arg GetLanguageVersionParams) (LanguageVersion, error) {
	row := q.db.QueryRow(ctx, getLanguageVersion, arg.LanguageID, arg.Version)
	var i LanguageVersion
	err := row.Scan(
		&i.ID,
		&i.LanguageID,
		&i.Version,
		&i.NixPackageName,
		&i.FlakeTemplate,
		&i.ScriptTemplate,
		&i.DefaultCode,
		&i.SearchQuery,
	)
	return i, err
}

const getLanguageVersionByID = `-- name: GetLanguageVersionByID :one
SELECT id, language_id, version, nix_package_name, flake_template, script_template, default_code, search_query  
FROM language_versions 
WHERE id = $1
`

func (q *Queries) GetLanguageVersionByID(ctx context.Context, id int64) (LanguageVersion, error) {
	row := q.db.QueryRow(ctx, getLanguageVersionByID, id)
	var i LanguageVersion
	err := row.Scan(
		&i.ID,
		&i.LanguageID,
		&i.Version,
		&i.NixPackageName,
		&i.FlakeTemplate,
		&i.ScriptTemplate,
		&i.DefaultCode,
		&i.SearchQuery,
	)
	return i, err
}

const getVersionsByLanguageID = `-- name: GetVersionsByLanguageID :many
SELECT id, language_id, version, nix_package_name, flake_template, script_template, default_code, search_query 
FROM language_versions 
WHERE language_id = $1
`

func (q *Queries) GetVersionsByLanguageID(ctx context.Context, languageID int64) ([]LanguageVersion, error) {
	rows, err := q.db.Query(ctx, getVersionsByLanguageID, languageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LanguageVersion
	for rows.Next() {
		var i LanguageVersion
		if err := rows.Scan(
			&i.ID,
			&i.LanguageID,
			&i.Version,
			&i.NixPackageName,
			&i.FlakeTemplate,
			&i.ScriptTemplate,
			&i.DefaultCode,
			&i.SearchQuery,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateLanguage = `-- name: UpdateLanguage :one
UPDATE languages 
SET name = $2, extension = $3, monaco_language = $4 
WHERE id = $1
returning id
`

type UpdateLanguageParams struct {
	ID             int64  `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	Extension      string `db:"extension" json:"extension"`
	MonacoLanguage string `db:"monaco_language" json:"monaco_language"`
}

func (q *Queries) UpdateLanguage(ctx context.Context, arg UpdateLanguageParams) (int64, error) {
	row := q.db.QueryRow(ctx, updateLanguage,
		arg.ID,
		arg.Name,
		arg.Extension,
		arg.MonacoLanguage,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateLanguageVersion = `-- name: UpdateLanguageVersion :one
UPDATE language_versions 
SET language_id = $2, nix_package_name = $3, flake_template = $4, script_template = $5, default_code = $6, search_query = $7, version = $8
WHERE id = $1 
returning id
`

type UpdateLanguageVersionParams struct {
	ID             int64  `db:"id" json:"id"`
	LanguageID     int64  `db:"language_id" json:"language_id"`
	NixPackageName string `db:"nix_package_name" json:"nix_package_name"`
	FlakeTemplate  string `db:"flake_template" json:"flake_template"`
	ScriptTemplate string `db:"script_template" json:"script_template"`
	DefaultCode    string `db:"default_code" json:"default_code"`
	SearchQuery    string `db:"search_query" json:"search_query"`
	Version        string `db:"version" json:"version"`
}

func (q *Queries) UpdateLanguageVersion(ctx context.Context, arg UpdateLanguageVersionParams) (int64, error) {
	row := q.db.QueryRow(ctx, updateLanguageVersion,
		arg.ID,
		arg.LanguageID,
		arg.NixPackageName,
		arg.FlakeTemplate,
		arg.ScriptTemplate,
		arg.DefaultCode,
		arg.SearchQuery,
		arg.Version,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

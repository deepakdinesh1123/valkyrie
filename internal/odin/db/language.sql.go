// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: language.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createLanguage = `-- name: CreateLanguage :one
INSERT INTO languages (name, extension, monaco_language, default_code) 
VALUES ($1, $2, $3, $4) 
RETURNING id
`

type CreateLanguageParams struct {
	Name           string `db:"name" json:"name"`
	Extension      string `db:"extension" json:"extension"`
	MonacoLanguage string `db:"monaco_language" json:"monaco_language"`
	DefaultCode    string `db:"default_code" json:"default_code"`
}

func (q *Queries) CreateLanguage(ctx context.Context, arg CreateLanguageParams) (int64, error) {
	row := q.db.QueryRow(ctx, createLanguage,
		arg.Name,
		arg.Extension,
		arg.MonacoLanguage,
		arg.DefaultCode,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createLanguageVersion = `-- name: CreateLanguageVersion :one
INSERT INTO language_versions (
    language_id, version, nix_package_name, template, search_query, default_version
) 
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING id
`

type CreateLanguageVersionParams struct {
	LanguageID     int64       `db:"language_id" json:"language_id"`
	Version        string      `db:"version" json:"version"`
	NixPackageName string      `db:"nix_package_name" json:"nix_package_name"`
	Template       pgtype.Text `db:"template" json:"template"`
	SearchQuery    string      `db:"search_query" json:"search_query"`
	DefaultVersion bool        `db:"default_version" json:"default_version"`
}

func (q *Queries) CreateLanguageVersion(ctx context.Context, arg CreateLanguageVersionParams) (int64, error) {
	row := q.db.QueryRow(ctx, createLanguageVersion,
		arg.LanguageID,
		arg.Version,
		arg.NixPackageName,
		arg.Template,
		arg.SearchQuery,
		arg.DefaultVersion,
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
SELECT id, language_id, version, nix_package_name, template, search_query, default_version
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
			&i.Template,
			&i.SearchQuery,
			&i.DefaultVersion,
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
SELECT id, name, extension, monaco_language, template, default_code
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
			&i.Template,
			&i.DefaultCode,
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

const getDefaultVersion = `-- name: GetDefaultVersion :one
SELECT id, language_id, version, nix_package_name, template, search_query, default_version FROM language_versions WHERE default_version = true AND language_id = $1
`

func (q *Queries) GetDefaultVersion(ctx context.Context, languageID int64) (LanguageVersion, error) {
	row := q.db.QueryRow(ctx, getDefaultVersion, languageID)
	var i LanguageVersion
	err := row.Scan(
		&i.ID,
		&i.LanguageID,
		&i.Version,
		&i.NixPackageName,
		&i.Template,
		&i.SearchQuery,
		&i.DefaultVersion,
	)
	return i, err
}

const getLanguageByID = `-- name: GetLanguageByID :one
SELECT id, name, extension, monaco_language, template, default_code
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
		&i.Template,
		&i.DefaultCode,
	)
	return i, err
}

const getLanguageByName = `-- name: GetLanguageByName :one
SELECT id, name, extension, monaco_language, template, default_code from languages WHERE name = $1
`

func (q *Queries) GetLanguageByName(ctx context.Context, name string) (Language, error) {
	row := q.db.QueryRow(ctx, getLanguageByName, name)
	var i Language
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Extension,
		&i.MonacoLanguage,
		&i.Template,
		&i.DefaultCode,
	)
	return i, err
}

const getLanguageVersion = `-- name: GetLanguageVersion :one
SELECT id, language_id, version, nix_package_name, template, search_query, default_version 
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
		&i.Template,
		&i.SearchQuery,
		&i.DefaultVersion,
	)
	return i, err
}

const getLanguageVersionByID = `-- name: GetLanguageVersionByID :one
SELECT id, language_id, version, nix_package_name, template, search_query, default_version 
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
		&i.Template,
		&i.SearchQuery,
		&i.DefaultVersion,
	)
	return i, err
}

const getVersionsByLanguageID = `-- name: GetVersionsByLanguageID :many
SELECT id, language_id, version, nix_package_name, template, search_query, default_version 
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
			&i.Template,
			&i.SearchQuery,
			&i.DefaultVersion,
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

type InsertLanguageVersionsParams struct {
	LanguageID     int64       `db:"language_id" json:"language_id"`
	Version        string      `db:"version" json:"version"`
	NixPackageName string      `db:"nix_package_name" json:"nix_package_name"`
	Template       pgtype.Text `db:"template" json:"template"`
	SearchQuery    string      `db:"search_query" json:"search_query"`
	DefaultVersion bool        `db:"default_version" json:"default_version"`
}

type InsertLanguagesParams struct {
	Name           string `db:"name" json:"name"`
	Extension      string `db:"extension" json:"extension"`
	MonacoLanguage string `db:"monaco_language" json:"monaco_language"`
	DefaultCode    string `db:"default_code" json:"default_code"`
	Template       string `db:"template" json:"template"`
}

const truncateLanguageVersions = `-- name: TruncateLanguageVersions :exec
TRUNCATE TABLE language_versions CASCADE
`

func (q *Queries) TruncateLanguageVersions(ctx context.Context) error {
	_, err := q.db.Exec(ctx, truncateLanguageVersions)
	return err
}

const truncateLanguages = `-- name: TruncateLanguages :exec
TRUNCATE TABLE languages CASCADE
`

func (q *Queries) TruncateLanguages(ctx context.Context) error {
	_, err := q.db.Exec(ctx, truncateLanguages)
	return err
}

const updateLanguage = `-- name: UpdateLanguage :one
UPDATE languages 
SET name = $2, extension = $3, monaco_language = $4, default_code = $5
WHERE id = $1
returning id
`

type UpdateLanguageParams struct {
	ID             int64  `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	Extension      string `db:"extension" json:"extension"`
	MonacoLanguage string `db:"monaco_language" json:"monaco_language"`
	DefaultCode    string `db:"default_code" json:"default_code"`
}

func (q *Queries) UpdateLanguage(ctx context.Context, arg UpdateLanguageParams) (int64, error) {
	row := q.db.QueryRow(ctx, updateLanguage,
		arg.ID,
		arg.Name,
		arg.Extension,
		arg.MonacoLanguage,
		arg.DefaultCode,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateLanguageVersion = `-- name: UpdateLanguageVersion :one
UPDATE language_versions 
SET language_id = $2, nix_package_name = $3, template = $4, search_query = $5, version = $6, default_version = $7
WHERE id = $1 
returning id
`

type UpdateLanguageVersionParams struct {
	ID             int64       `db:"id" json:"id"`
	LanguageID     int64       `db:"language_id" json:"language_id"`
	NixPackageName string      `db:"nix_package_name" json:"nix_package_name"`
	Template       pgtype.Text `db:"template" json:"template"`
	SearchQuery    string      `db:"search_query" json:"search_query"`
	Version        string      `db:"version" json:"version"`
	DefaultVersion bool        `db:"default_version" json:"default_version"`
}

func (q *Queries) UpdateLanguageVersion(ctx context.Context, arg UpdateLanguageVersionParams) (int64, error) {
	row := q.db.QueryRow(ctx, updateLanguageVersion,
		arg.ID,
		arg.LanguageID,
		arg.NixPackageName,
		arg.Template,
		arg.SearchQuery,
		arg.Version,
		arg.DefaultVersion,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: package.sql

package db

import (
	"context"
)

const packagesExist = `-- name: PackagesExist :one
WITH existing_packages AS (
    SELECT name
    FROM packages
    WHERE language = $2::text
      AND name = ANY($1::text[])
)
SELECT 
    COUNT(*) = array_length($1::text[], 1) AS exists,
    ARRAY(
        SELECT unnest($1::text[]) 
        EXCEPT 
        SELECT name FROM existing_packages
    )::text[] AS nonexisting_packages
FROM existing_packages
`

type PackagesExistParams struct {
	Packages []string `db:"packages" json:"packages"`
	Language string   `db:"language" json:"language"`
}

type PackagesExistRow struct {
	Exists              bool     `db:"exists" json:"exists"`
	NonexistingPackages []string `db:"nonexisting_packages" json:"nonexisting_packages"`
}

func (q *Queries) PackagesExist(ctx context.Context, arg PackagesExistParams) (PackagesExistRow, error) {
	row := q.db.QueryRow(ctx, packagesExist, arg.Packages, arg.Language)
	var i PackagesExistRow
	err := row.Scan(&i.Exists, &i.NonexistingPackages)
	return i, err
}

const searchLanguagePackages = `-- name: SearchLanguagePackages :many
SELECT
    name,
    version
FROM packages
WHERE
    language = $1::text  
    AND (
        tsv_search @@ plainto_tsquery('english', $2::text) OR  -- Full-text search
        name ILIKE '%' || $2::text || '%'                       -- Contains search
    )
ORDER BY name ASC
LIMIT 50
`

type SearchLanguagePackagesParams struct {
	Language    string `db:"language" json:"language"`
	Searchquery string `db:"searchquery" json:"searchquery"`
}

type SearchLanguagePackagesRow struct {
	Name    string `db:"name" json:"name"`
	Version string `db:"version" json:"version"`
}

func (q *Queries) SearchLanguagePackages(ctx context.Context, arg SearchLanguagePackagesParams) ([]SearchLanguagePackagesRow, error) {
	rows, err := q.db.Query(ctx, searchLanguagePackages, arg.Language, arg.Searchquery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchLanguagePackagesRow
	for rows.Next() {
		var i SearchLanguagePackagesRow
		if err := rows.Scan(&i.Name, &i.Version); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchSystemPackages = `-- name: SearchSystemPackages :many
SELECT
    name,
    version
FROM packages
WHERE
    pkgType = 'system'
    AND (
        tsv_search @@ plainto_tsquery('english', $1) OR  -- Full-text search
        name ILIKE '%' || $1 || '%'                       -- Contains search
    )
ORDER BY name ASC
LIMIT 50
`

type SearchSystemPackagesRow struct {
	Name    string `db:"name" json:"name"`
	Version string `db:"version" json:"version"`
}

func (q *Queries) SearchSystemPackages(ctx context.Context, plaintoTsquery string) ([]SearchSystemPackagesRow, error) {
	rows, err := q.db.Query(ctx, searchSystemPackages, plaintoTsquery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchSystemPackagesRow
	for rows.Next() {
		var i SearchSystemPackagesRow
		if err := rows.Scan(&i.Name, &i.Version); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

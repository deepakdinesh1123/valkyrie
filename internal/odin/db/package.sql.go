// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: package.sql

package db

import (
	"context"
)

const searchLanguagePackages = `-- name: SearchLanguagePackages :many
SELECT
    name,
    version
FROM packages
WHERE
    language = $1::text  
    AND tsv_search @@ plainto_tsquery('english', $2::text)  
ORDER BY name ASC
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
    AND tsv_search @@ plainto_tsquery('english', $1)
ORDER BY name ASC
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
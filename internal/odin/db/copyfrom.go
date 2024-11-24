// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: copyfrom.go

package db

import (
	"context"
)

// iteratorForInsertLanguageVersions implements pgx.CopyFromSource.
type iteratorForInsertLanguageVersions struct {
	rows                 []InsertLanguageVersionsParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertLanguageVersions) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertLanguageVersions) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].LanguageID,
		r.rows[0].Version,
		r.rows[0].NixPackageName,
		r.rows[0].Template,
		r.rows[0].SearchQuery,
		r.rows[0].DefaultVersion,
	}, nil
}

func (r iteratorForInsertLanguageVersions) Err() error {
	return nil
}

func (q *Queries) InsertLanguageVersions(ctx context.Context, arg []InsertLanguageVersionsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"language_versions"}, []string{"language_id", "version", "nix_package_name", "template", "search_query", "default_version"}, &iteratorForInsertLanguageVersions{rows: arg})
}

// iteratorForInsertLanguages implements pgx.CopyFromSource.
type iteratorForInsertLanguages struct {
	rows                 []InsertLanguagesParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertLanguages) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertLanguages) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Extension,
		r.rows[0].MonacoLanguage,
		r.rows[0].DefaultCode,
		r.rows[0].Template,
	}, nil
}

func (r iteratorForInsertLanguages) Err() error {
	return nil
}

func (q *Queries) InsertLanguages(ctx context.Context, arg []InsertLanguagesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"languages"}, []string{"name", "extension", "monaco_language", "default_code", "template"}, &iteratorForInsertLanguages{rows: arg})
}

// iteratorForInsertPackages implements pgx.CopyFromSource.
type iteratorForInsertPackages struct {
	rows                 []InsertPackagesParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertPackages) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertPackages) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Version,
		r.rows[0].Pkgtype,
		r.rows[0].Language,
		r.rows[0].StorePath,
	}, nil
}

func (r iteratorForInsertPackages) Err() error {
	return nil
}

func (q *Queries) InsertPackages(ctx context.Context, arg []InsertPackagesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"packages"}, []string{"name", "version", "pkgtype", "language", "store_path"}, &iteratorForInsertPackages{rows: arg})
}

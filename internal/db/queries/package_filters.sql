-- name: TruncateSysPkgFilters :exec
TRUNCATE TABLE system_package_filters CASCADE;

-- name: GetSysPkgFilters :many
SELECT * from system_package_filters;

-- name: BulkInsertSystemPackageFilters :copyfrom
INSERT INTO system_package_filters (
    filter_type,
    package_string
) VALUES (
    $1,
    $2
);

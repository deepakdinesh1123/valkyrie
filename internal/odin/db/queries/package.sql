create table if not exists packages (
    package_id bigint primary key default nextval('packages_id_seq'),
    name text not null,
    version text not null,
    pkgType text not null,
    language text,
    store_path text,
    tsv_search TSVECTOR
);

-- name: SearchSystemPackages :many
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
LIMIT 50;

-- name: SearchLanguagePackages :many
SELECT
    name,
    version
FROM packages
WHERE
    language = @Language::text  
    AND (
        tsv_search @@ plainto_tsquery('english', @SearchQuery::text) OR  -- Full-text search
        name ILIKE '%' || @SearchQuery::text || '%'                       -- Contains search
    )
ORDER BY name ASC
LIMIT 50;


-- name: PackagesExist :one
WITH existing_packages AS (
    SELECT name
    FROM packages
    WHERE language = @language::text
      AND name = ANY(@packages::text[])
)
SELECT 
    COUNT(*) = array_length(@packages::text[], 1) AS exists,
    ARRAY(
        SELECT unnest(@packages::text[]) 
        EXCEPT 
        SELECT name FROM existing_packages
    )::text[] AS nonexisting_packages
FROM existing_packages;

-- name: GetPackageStorePaths :many
SELECT name, store_path
FROM packages
WHERE name = ANY(@packages::text[]);



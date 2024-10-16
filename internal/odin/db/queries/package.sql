create table packages (
    package_id bigint primary key default nextval('packages_id_seq'),
    name text not null,
    version text not null,
    pkgType text not null,
    language text,
    tsv_search TSVECTOR
);

-- name: SearchSystemPackages :many
SELECT
    name,
    version
FROM packages
WHERE
    pkgType = 'system'
    AND tsv_search @@ plainto_tsquery('english', $1)
ORDER BY name ASC;


-- name: SearchLanguagePackages :many
SELECT
    name,
    version
FROM packages
WHERE
    language = @Language::text  
    AND tsv_search @@ plainto_tsquery('english', @SearchQuery::text)  
ORDER BY name ASC;

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



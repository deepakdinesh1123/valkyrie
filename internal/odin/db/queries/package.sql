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
    tsv_search @@ plainto_tsquery('english', $1)  
    AND tsv_search @@ plainto_tsquery('english', $2)  
ORDER BY name ASC;







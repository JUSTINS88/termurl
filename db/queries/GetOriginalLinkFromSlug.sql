-- name: GetOriginalLinkFromSlug :one
SELECT original_link FROM links
WHERE slug = $1;
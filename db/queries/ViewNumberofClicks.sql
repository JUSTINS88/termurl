-- name: GetNumberOfClicks :one
SELECT original_link, slug, clicks, created_at FROM links
WHERE slug = $1;
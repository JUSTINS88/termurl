-- name: UpdateClicks :execrows
UPDATE links
SET clicks = clicks + 1
WHERE SLUG = $1;
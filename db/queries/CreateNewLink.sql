-- name: CreateNewLink :one
INSERT INTO links (id, original_link, slug, clicks, machine_key, created_at)
VALUES (
    gen_random_uuid(),
        $1,
        $2,
        0,
        $3,
        NOW()
)
RETURNING *;
-- name: CheckMachineKey :one
-- name: CheckMachineKey :one
INSERT INTO users (machine_key, created_at)
VALUES ($1, NOW())
ON CONFLICT (machine_key) DO UPDATE SET machine_key = EXCLUDED.machine_key
RETURNING machine_key;

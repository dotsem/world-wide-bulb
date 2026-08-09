-- name: InsertToggle :one
INSERT INTO toggles (state, reason, ip_hash)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetRecentToggles :many
SELECT * FROM toggles
ORDER BY created_at DESC
LIMIT ?;

-- name: GetLatestToggle :one
SELECT * FROM toggles
ORDER BY created_at DESC
LIMIT 1;

-- name: PruneOldToggles :exec
DELETE FROM toggles
WHERE id NOT IN (
    SELECT id FROM toggles ORDER BY id DESC LIMIT ?
);

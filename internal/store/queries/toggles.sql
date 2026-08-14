-- name: InsertToggle :one
INSERT INTO toggles (uuid, state, ip_hash)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateToggleReason :execresult
UPDATE toggles
SET reason = ?
WHERE uuid = ?;

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

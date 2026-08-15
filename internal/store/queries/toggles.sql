-- name: InsertToggle :one
INSERT INTO toggles (uuid, state, ip_hash)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateToggleReason :execresult
UPDATE toggles
SET reason = ?
WHERE uuid = ? AND (reason IS NULL OR reason = '');

-- name: GetToggleByUUID :one
SELECT * FROM toggles
WHERE uuid = ?;

-- name: GetRecentToggles :many
SELECT * FROM toggles
ORDER BY created_at DESC
LIMIT ?;

-- name: GetTogglesBefore :many
SELECT * FROM toggles
WHERE id < ?
ORDER BY id DESC
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

-- name: GetEntity :one
SELECT id, field1, field2
FROM entities
WHERE id = $1;

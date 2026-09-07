-- name: CreateEntity :one
INSERT INTO entities (field1, field2)
VALUES ($1, $2)
RETURNING id, field1, field2;

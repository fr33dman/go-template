-- +goose Up
CREATE TABLE entities (
    id BIGSERIAL PRIMARY KEY,
    field1 TEXT NOT NULL,
    field2 INTEGER NOT NULL
);

-- +goose Down
DROP TABLE entities;

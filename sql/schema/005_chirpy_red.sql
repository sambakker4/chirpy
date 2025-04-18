-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red bool DEFAULT false;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_chirpy_red;

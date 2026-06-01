-- +goose Up
ALTER TABLE collections
ADD COLUMN "group" TEXT;

-- +goose Down
ALTER TABLE collections
DROP COLUMN "group";

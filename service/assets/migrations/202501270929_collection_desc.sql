-- +goose Up
ALTER TABLE collections
ADD description text;

-- +goose Down
ALTER TABLE collections
DROP COLUMN description;

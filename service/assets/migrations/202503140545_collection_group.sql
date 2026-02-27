-- +goose Up
ALTER TABLE collections
ADD group_name varchar(40);
UPDATE collections SET group_name = 'foxstream';

-- +goose Down
ALTER TABLE collections
DROP COLUMN group_name;

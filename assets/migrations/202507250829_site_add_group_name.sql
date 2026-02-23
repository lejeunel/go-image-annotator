-- +goose Up
ALTER TABLE sites
ADD group_name varchar(36) DEFAULT 'foxstream';

-- +goose Down
ALTER TABLE sites
DROP COLUMN group_name;

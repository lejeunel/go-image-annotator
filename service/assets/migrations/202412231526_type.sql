-- +goose Up
ALTER TABLE images
ADD image_type varchar(15);

-- +goose Down
ALTER TABLE images
DROP COLUMN image_type;

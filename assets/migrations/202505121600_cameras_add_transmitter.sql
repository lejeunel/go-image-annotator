-- +goose Up
ALTER TABLE cameras
ADD transmitter varchar(36);

-- +goose Down
ALTER TABLE cameras
DROP COLUMN transmitter;

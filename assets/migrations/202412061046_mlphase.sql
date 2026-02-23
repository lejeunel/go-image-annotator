-- +goose Up
CREATE TABLE IF NOT EXISTS mlphase (
    id varchar(36) NOT NULL,
    image_id varchar(36),
    collection_id varchar(36),
    phase varchar(10),
    PRIMARY KEY(id),
    FOREIGN KEY (image_id) REFERENCES images(id)
);

-- +goose Down
DROP TABLE mlphase;

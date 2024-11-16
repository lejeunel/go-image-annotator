-- +goose Up

CREATE TABLE IF NOT EXISTS images (
    id varchar(16),
    uri text,
    created_at text,
    updated_at text,
    sha256 varchar(64),
    width int,
    height int,
    mimetype varchar(40),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS labels (
    id varchar(16),
    name text,
    description text,
    created_at text,
    updated_at text,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS image_label_assoc (
    id varchar(16) PRIMARY KEY,
    image_id varchar(16),
    label_id varchar(16),
    created_at text,
    FOREIGN KEY (image_id) REFERENCES images(id),
    FOREIGN KEY (label_id) REFERENCES labels(id)
);

-- +goose Down

DROP TABLE images;
DROP TABLE labels;
DROP TABLE image_label_assoc;

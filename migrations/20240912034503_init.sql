-- +goose Up

CREATE TABLE IF NOT EXISTS images (
    id varchar(16),
    uri text,
    created_at text,
    updated_at text,
    sha256 varchar(64),
    width int,
    height int,
    mimetype varchar(40)
);

CREATE TABLE IF NOT EXISTS labels (
    id varchar(16),
    name text,
    description text,
    created_at text,
    updated_at text
);


-- +goose Down

DROP TABLE images;
DROP TABLE labels;

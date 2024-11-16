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



-- +goose Down

DROP TABLE images;
DROP TABLE sequences;

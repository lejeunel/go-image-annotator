-- +goose Up

CREATE TABLE IF NOT EXISTS images (
    id varchar(16),
    uri text,
    created_at varchar(30),
    updated_at varchar(30),
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
    created_at varchar(30),
    updated_at varchar(30),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS annotations (
    id varchar(16) PRIMARY KEY,
    image_id varchar(16),
    label_id varchar(16),
    author_email varchar(40),
    created_at varchar(30),
    updated_at varchar(30),
    shape_type varchar(16),
    shape_data varchar(200),
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    FOREIGN KEY (label_id) REFERENCES labels(id)
);


CREATE TABLE IF NOT EXISTS imagesets (
    id varchar(16) PRIMARY KEY,
    name text,
    created_at varchar(30),
    updated_at varchar(30)
);

CREATE TABLE IF NOT EXISTS image_set_assoc (
    id varchar(16) PRIMARY KEY,
    image_id varchar(16),
    set_id varchar(16),
    created_at varchar(30),
    FOREIGN KEY (image_id) REFERENCES images(id),
    FOREIGN KEY (set_id) REFERENCES imagesets(id)
);

-- +goose Down

DROP TABLE images;
DROP TABLE labels;
DROP TABLE annotations;
DROP TABLE polygons;
DROP TABLE imagesets;
DROP TABLE image_set_assoc;

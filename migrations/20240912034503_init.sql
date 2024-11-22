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
    author_email varchar(40),
    created_at text,
    FOREIGN KEY (image_id) REFERENCES images(id),
    FOREIGN KEY (label_id) REFERENCES labels(id)
);

CREATE TABLE IF NOT EXISTS polygons (
    id varchar(16) PRIMARY KEY,
    image_id varchar(16),
    label_id varchar(16),
    created_at text,
    updated_at text,
    type_ varchar(16),
    min_x int,
    min_y int,
    max_x int,
    max_y int,
    points text,
    FOREIGN KEY (image_id) REFERENCES images(id),
    FOREIGN KEY (label_id) REFERENCES labels(id)
);

CREATE TABLE IF NOT EXISTS imagesets (
    id varchar(16) PRIMARY KEY,
    name text,
    created_at text,
    updated_at text
);

CREATE TABLE IF NOT EXISTS image_set_assoc (
    id varchar(16) PRIMARY KEY,
    image_id varchar(16),
    set_id varchar(16),
    created_at text,
    FOREIGN KEY (image_id) REFERENCES images(id),
    FOREIGN KEY (set_id) REFERENCES imagesets(id)
);

-- +goose Down

DROP TABLE images;
DROP TABLE labels;
DROP TABLE image_label_assoc;
DROP TABLE polygons;
DROP TABLE imagesets;
DROP TABLE image_set_assoc;

-- +goose Up

CREATE TABLE IF NOT EXISTS images (
    id varchar(36),
    filename varchar(120),
    site_id varchar(36),
    camera_id varchar(36),
    created_at timestamp,
    updated_at timestamp,
    captured_at timestamp,
    sha256 varchar(64),
    width int,
    height int,
    mimetype varchar(40),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS labels (
    id varchar(36),
    parent_id varchar(36),
    name varchar(30) UNIQUE,
    description text,
    created_at timestamp,
    updated_at timestamp,
    FOREIGN KEY (parent_id) REFERENCES labels(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS collections (
    id varchar(36) PRIMARY KEY,
    name text UNIQUE,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS annotations (
    id varchar(36) PRIMARY KEY,
    image_id varchar(36),
    label_id varchar(36),
    collection_id varchar(36),
    author_email varchar(40),
    created_at timestamp,
    updated_at timestamp,
    shape_type varchar(16),
    shape_data varchar(200),
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    FOREIGN KEY (label_id) REFERENCES labels(id) ON DELETE RESTRICT,
    FOREIGN KEY (collection_id) REFERENCES collections(id)
);



CREATE TABLE IF NOT EXISTS image_collection_assoc (
    id varchar(36) PRIMARY KEY,
    image_id varchar(36),
    collection_id varchar(36),
    created_at timestamp,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    FOREIGN KEY (collection_id) REFERENCES collections(id)
);

CREATE TABLE IF NOT EXISTS sites (
    id varchar(36) PRIMARY KEY,
    name varchar(30) UNIQUE,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS cameras (
    id varchar(36) PRIMARY KEY,
    name varchar(30),
    site_id varchar(36),
    created_at timestamp,
    updated_at timestamp,
    FOREIGN KEY (site_id) REFERENCES sites(id)
);

-- +goose Down

DROP TABLE images;
DROP TABLE labels;
DROP TABLE annotations;
DROP TABLE collections;
DROP TABLE image_collection_assoc;
DROP TABLE sites;
DROP TABLE cameras;

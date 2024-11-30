-- +goose Up

CREATE TABLE IF NOT EXISTS images (
    id varchar(36),
    uri varchar(60),
    created_at timestamp,
    updated_at timestamp,
    sha256 varchar(64),
    width int,
    height int,
    mimetype varchar(40),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS labels (
    id varchar(36),
    name varchar(30) UNIQUE,
    description text,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY (id)
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
    FOREIGN KEY (label_id) REFERENCES labels(id),
    FOREIGN KEY (collection_id) REFERENCES collections(id)
);


CREATE TABLE IF NOT EXISTS collections (
    id varchar(36) PRIMARY KEY,
    name text UNIQUE,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS image_collection_assoc (
    id varchar(36) PRIMARY KEY,
    image_id varchar(36),
    collection_id varchar(36),
    created_at timestamp,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    FOREIGN KEY (collection_id) REFERENCES collections(id)
);

-- +goose Down

DROP TABLE images;
DROP TABLE labels;
DROP TABLE annotations;
DROP TABLE collections;
DROP TABLE image_collection_assoc;

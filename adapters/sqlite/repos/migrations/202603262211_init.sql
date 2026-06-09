-- +goose Up

CREATE TABLE IF NOT EXISTS labels (
    id varchar(36),
    name varchar(30) not null unique,
    description text,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_labels_name ON labels(name);

CREATE TABLE IF NOT EXISTS collections (
    id varchar(36),
    name varchar(30) not null unique,
    description text,
    created_at DATETIME,
    group_id varchar(36) NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_collections_name ON collections(name);

CREATE TABLE IF NOT EXISTS groups (
    id varchar(36),
    name varchar(30) not null unique,
    description text,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_groups_name ON groups(name);

CREATE TABLE IF NOT EXISTS images (
    id varchar(36),
    hash varchar(128),
    mimetype TEXT,
    width INTEGER,
    height INTEGER,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_images_hash ON images(hash);

CREATE TABLE IF NOT EXISTS images_collections (
  image_id varchar(36) REFERENCES images(id),
  collection_id varchar(36) REFERENCES collections(id),
  PRIMARY KEY (image_id, collection_id)
);

CREATE TABLE IF NOT EXISTS annotations (
  id varchar(36),
  image_id varchar(36) REFERENCES images(id),
  collection_id varchar(36) REFERENCES collections(id),
  label_id varchar(36) REFERENCES labels(id),
  type varchar(15),
  coordinates varchar(100),
  PRIMARY KEY (id)
);

CREATE INDEX idx_annotations_image_collection ON annotations(image_id,collection_id);

CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BLOB NOT NULL,
	expiry REAL NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions(expiry);

CREATE TABLE users (
	id varchar(60) PRIMARY KEY,
    roles TEXT
);

CREATE TABLE IF NOT EXISTS users_groups (
  user_id varchar(36) REFERENCES users(id),
  group_id varchar(36) REFERENCES groups(id),
  PRIMARY KEY (user_id, group_id)
);
-- +goose Down

DROP TABLE labels;
DROP TABLE collections;
DROP TABLE images_collections;
DROP TABLE images;
DROP TABLE annotations;
DROP TABLE sessions;
DROP TABLE groups;

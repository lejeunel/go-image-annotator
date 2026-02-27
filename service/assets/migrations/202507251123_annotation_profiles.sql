-- +goose Up

CREATE TABLE IF NOT EXISTS annotation_profiles (
    id varchar(36),
    name varchar(100) UNIQUE,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS annotation_profile_label_assoc (
    id varchar(36) PRIMARY KEY,
    profile_id varchar(36),
    label_id varchar(36),
    FOREIGN KEY (profile_id) REFERENCES annotation_profiles(id) ON DELETE CASCADE,
    FOREIGN KEY (label_id) REFERENCES labels(id)
);

-- +goose Down

DROP TABLE annotation_profiles;
DROP TABLE annotation_profile_label_assoc;

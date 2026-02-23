-- +goose Up
CREATE INDEX idx_image_site ON images(site_id);
CREATE INDEX idx_image_camera ON images(camera_id);

CREATE INDEX idx_label_parent_label ON labels(parent_id);

CREATE INDEX idx_annotation_image ON annotations(image_id);
CREATE INDEX idx_annotation_collection ON annotations(collection_id);

CREATE INDEX idx_image_collection_image ON image_collection_assoc(image_id);
CREATE INDEX idx_image_collection_collection ON image_collection_assoc(collection_id);

CREATE INDEX idx_mlphase_image ON mlphase(image_id);
CREATE INDEX idx_mlphase_collection ON mlphase(collection_id);

-- +goose Down
DROP INDEX idx_image_site;
DROP INDEX idx_image_camera;

DROP INDEX idx_label_parent_label;

DROP INDEX idx_annotation_image;
DROP INDEX idx_annotation_collection;

DROP INDEX idx_image_collection_image;
DROP INDEX idx_image_collection_collection;

DROP INDEX idx_mlphase_image;
DROP INDEX idx_mlphase_collection;

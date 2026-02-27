-- +goose Up
UPDATE cameras
SET transmitter = '';

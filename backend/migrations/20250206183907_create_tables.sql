-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS containers (
    id SERIAL PRIMARY KEY,
    container_id VARCHAR(64) NOT NULL,
    container_name VARCHAR(256),
    ip_address VARCHAR(64),
    status VARCHAR(32),
    updated_at TIMESTAMP NOT NULL
);

ALTER TABLE containers 
  ADD CONSTRAINT containers_container_id_key UNIQUE (container_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS containers;
-- +goose StatementEnd

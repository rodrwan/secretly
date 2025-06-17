-- +goose Up
-- +goose StatementBegin
CREATE TABLE environment_values (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    environment_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (environment_id) REFERENCES environment (id)
);

CREATE INDEX idx_environment_values_environment_id ON environment_values (environment_id);
CREATE INDEX idx_environment_values_key ON environment_values (key);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE environment_values;
-- +goose StatementEnd

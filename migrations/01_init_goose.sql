-- +goose Up
CREATE TABLE IF NOT EXISTS goose_db_version (
                                                id INTEGER PRIMARY KEY AUTOINCREMENT,
                                                version_id INTEGER NOT NULL,
                                                is_applied INTEGER NOT NULL,  -- SQLite не имеет типа BOOLEAN
                                                tstamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS goose_db_version_version_id ON goose_db_version(version_id);

-- +goose Down
DROP TABLE goose_db_version;
CREATE TABLE IF NOT EXISTS toggles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    state BOOLEAN NOT NULL,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_hash TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_toggles_created_at ON toggles(created_at DESC);
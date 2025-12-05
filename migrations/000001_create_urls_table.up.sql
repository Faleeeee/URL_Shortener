CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    alias VARCHAR(16) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    click_count BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE UNIQUE INDEX idx_alias ON urls(alias);
CREATE INDEX idx_created_at ON urls(created_at);

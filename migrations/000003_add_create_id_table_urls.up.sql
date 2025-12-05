ALTER TABLE urls
ADD COLUMN create_id BIGINT NOT NULL;

-- Indexes for performance
CREATE INDEX idx_create_id ON urls(create_id);
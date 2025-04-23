CREATE TABLE IF NOT EXISTS campaign_metrics (
    id SERIAL PRIMARY KEY,
    campaign_id TEXT NOT NULL,
    platform TEXT NOT NULL,
    impressions INT DEFAULT 0,
    clicks INT DEFAULT 0,
    conversions INT DEFAULT 0,
    cost NUMERIC(10, 2) DEFAULT 0.00,
    revenue NUMERIC(10, 2) DEFAULT 0.00,
    timestamp TIMESTAMP NOT NULL
);
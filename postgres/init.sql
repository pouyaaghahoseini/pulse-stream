CREATE TABLE IF NOT EXISTS platform_analytics (
    platform TEXT PRIMARY KEY,
    total_posts INT NOT NULL,
    total_engagement INT NOT NULL,
    average_score DOUBLE PRECISION NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
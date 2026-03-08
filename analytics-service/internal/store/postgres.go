package store

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type PlatformStats struct {
	Platform        string
	TotalPosts      int
	TotalEngagement int
	AverageScore    float64
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connectionString string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) UpsertPlatformStats(stats PlatformStats) error {
	query := `
		INSERT INTO platform_analytics (
			platform,
			total_posts,
			total_engagement,
			average_score,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (platform)
		DO UPDATE SET
			total_posts = EXCLUDED.total_posts,
			total_engagement = EXCLUDED.total_engagement,
			average_score = EXCLUDED.average_score,
			updated_at = EXCLUDED.updated_at;
	`

	_, err := s.db.Exec(
		query,
		stats.Platform,
		stats.TotalPosts,
		stats.TotalEngagement,
		stats.AverageScore,
		time.Now().UTC(),
	)

	return err
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}
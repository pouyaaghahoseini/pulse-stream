package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PlatformAnalytics struct {
	Platform        string  `json:"platform"`
	TotalPosts      int     `json:"total_posts"`
	TotalEngagement int     `json:"total_engagement"`
	AverageScore    float64 `json:"average_score"`
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

func (s *PostgresStore) GetAllPlatformAnalytics() ([]PlatformAnalytics, error) {
	query := `
		SELECT platform, total_posts, total_engagement, average_score
		FROM platform_analytics
		ORDER BY platform;
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PlatformAnalytics

	for rows.Next() {
		var item PlatformAnalytics

		err := rows.Scan(
			&item.Platform,
			&item.TotalPosts,
			&item.TotalEngagement,
			&item.AverageScore,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}
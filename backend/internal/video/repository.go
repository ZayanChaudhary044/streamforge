package video

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(ctx context.Context, title string) (*Video, error) {
	id := uuid.New().String()

	const query = `
INSERT INTO videos (id, title)
VALUES ($1, $2)
RETURNING created_at;
`
	var createdAt time.Time

	if err := r.DB.QueryRowContext(ctx, query, id, title).Scan(&createdAt); err != nil {
		return nil, err
	}

	return &Video{
		ID:        id,
		Title:     title,
		CreatedAt: createdAt,
	}, nil
}

func (r *Repository) List(ctx context.Context) ([]Video, error) {
	const query = `
SELECT id, title, created_at
FROM videos
ORDER BY created_at DESC;
`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos := []Video{}
	for rows.Next() {
		var v Video
		if err := rows.Scan(&v.ID, &v.Title, &v.CreatedAt); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return videos, nil
}

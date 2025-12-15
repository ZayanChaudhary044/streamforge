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

// Create a placeholder "processing" video row
func (r *Repository) CreateProcessing(ctx context.Context, id, title string) (*Video, error) {
	const q = `
INSERT INTO videos (id, title, status)
VALUES ($1, $2, 'processing')
RETURNING created_at, status;
`
	var createdAt time.Time
	var status string

	if err := r.DB.QueryRowContext(ctx, q, id, title).Scan(&createdAt, &status); err != nil {
		return nil, err
	}

	return &Video{
		ID:        id,
		Title:     title,
		Status:    status,
		CreatedAt: createdAt,
	}, nil
}

func (r *Repository) EnqueueJob(ctx context.Context, jobID, videoID, rawPath string) error {
	const q = `
INSERT INTO video_jobs (id, video_id, raw_path, status)
VALUES ($1, $2, $3, 'queued');
`
	_, err := r.DB.ExecContext(ctx, q, jobID, videoID, rawPath)
	return err
}

func (r *Repository) MarkVideoReady(ctx context.Context, videoID, filename, thumbnail string) error {
	const q = `
UPDATE videos
SET filename=$2, thumbnail=$3, status='ready', error=NULL
WHERE id=$1;
`
	_, err := r.DB.ExecContext(ctx, q, videoID, filename, thumbnail)
	return err
}

func (r *Repository) MarkVideoFailed(ctx context.Context, videoID, msg string) error {
	const q = `
UPDATE videos
SET status='failed', error=$2
WHERE id=$1;
`
	_, err := r.DB.ExecContext(ctx, q, videoID, msg)
	return err
}

func (r *Repository) List(ctx context.Context) ([]Video, error) {
	const q = `
SELECT id, title, filename, thumbnail, status, error, created_at
FROM videos
ORDER BY created_at DESC;
`
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video

	for rows.Next() {
		var v Video
		var filename sql.NullString
		var thumb sql.NullString
		var errMsg sql.NullString

		if err := rows.Scan(
			&v.ID,
			&v.Title,
			&filename,
			&thumb,
			&v.Status,
			&errMsg,
			&v.CreatedAt,
		); err != nil {
			return nil, err
		}

		if filename.Valid {
			v.Filename = &filename.String
		}
		if thumb.Valid {
			v.Thumbnail = &thumb.String
		}
		if errMsg.Valid {
			v.Error = &errMsg.String
		}

		videos = append(videos, v)
	}

	return videos, nil
}

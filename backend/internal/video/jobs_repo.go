package video

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

func (r *Repository) ClaimNextJob(ctx context.Context) (*Job, error) {
	// Atomically:
	// - select the oldest queued job (skip locked)
	// - mark it running
	// - increment attempts
	const q = `
WITH cte AS (
  SELECT id
  FROM video_jobs
  WHERE status='queued'
  ORDER BY created_at
  FOR UPDATE SKIP LOCKED
  LIMIT 1
)
UPDATE video_jobs j
SET status='running',
    attempts = attempts + 1,
    updated_at = now()
FROM cte
WHERE j.id = cte.id
RETURNING j.id, j.video_id, j.raw_path, j.attempts;
`

	var job Job
	err := r.DB.QueryRowContext(ctx, q).Scan(&job.ID, &job.VideoID, &job.RawPath, &job.Attempts)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // nothing to do
		}
		return nil, err
	}

	return &job, nil
}

func (r *Repository) MarkJobDone(ctx context.Context, jobID string) error {
	const q = `
UPDATE video_jobs
SET status='done', updated_at=now(), last_error=NULL
WHERE id=$1;
`
	_, err := r.DB.ExecContext(ctx, q, jobID)
	return err
}

func (r *Repository) MarkJobFailed(ctx context.Context, jobID string, msg string) error {
	const q = `
UPDATE video_jobs
SET status='failed', last_error=$2, updated_at=now()
WHERE id=$1;
`
	_, err := r.DB.ExecContext(ctx, q, jobID, msg)
	return err
}

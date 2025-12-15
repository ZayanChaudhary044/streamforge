package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ZayanChaudhary044/streamforge/backend/internal/database"
	"github.com/ZayanChaudhary044/streamforge/backend/internal/video"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	db := database.Open(dsn)
	defer db.Close()

	repo := video.NewRepository(db)

	// Ensure output directories exist
	if err := os.MkdirAll(filepath.Join(video.UploadDir, "thumbnails"), 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(video.UploadDir, "raw"), 0755); err != nil {
		log.Fatal(err)
	}

	log.Println("worker listening for jobs...")

	for {
		ctx := context.Background()

		job, err := repo.ClaimNextJob(ctx)
		if err != nil {
			log.Printf("claim job error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if job == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("picked job=%s video=%s raw=%s attempts=%d", job.ID, job.VideoID, job.RawPath, job.Attempts)

		if err := processJob(ctx, repo, job); err != nil {
			log.Printf("job failed job=%s: %v", job.ID, err)

			_ = repo.MarkVideoFailed(ctx, job.VideoID, err.Error())
			_ = repo.MarkJobFailed(ctx, job.ID, err.Error())
			continue
		}

		_ = repo.MarkJobDone(ctx, job.ID)
		log.Printf("job done job=%s video=%s", job.ID, job.VideoID)
	}
}

func processJob(ctx context.Context, repo *video.Repository, job *video.Job) error {
	finalFilename := job.VideoID + ".mp4"
	finalPath := filepath.Join(video.UploadDir, finalFilename)

	thumbName := job.VideoID + ".jpg"
	thumbPath := filepath.Join(video.UploadDir, "thumbnails", thumbName)

	// 1) Transcode raw -> mp4
	transcodeCmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", job.RawPath,
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "128k",
		finalPath,
	)
	transcodeCmd.Stdout = os.Stdout
	transcodeCmd.Stderr = os.Stderr
	if err := transcodeCmd.Run(); err != nil {
		return err
	}

	// 2) Thumbnail
	thumbCmd := exec.Command(
	"ffmpeg",
	"-y",
	"-ss", "00:00:01",
	"-i", finalPath,
	"-frames:v", "1",
	thumbPath,
)

	thumbCmd.Stdout = os.Stdout
	thumbCmd.Stderr = os.Stderr
	if err := thumbCmd.Run(); err != nil {
		return err
	}

	// 3) Update video row
	if err := repo.MarkVideoReady(ctx, job.VideoID, finalFilename, thumbName); err != nil {
		return err
	}

	// Optional: remove raw input after success
	_ = os.Remove(job.RawPath)

	return nil
}

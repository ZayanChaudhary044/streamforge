package video

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

const UploadDir = "./uploads"

type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{Repo: repo}
}

// List videos
func (h *Handler) ListVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := h.Repo.List(r.Context())
	if err != nil {
		http.Error(w, "failed to list videos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// Upload video + generate thumbnail
func (h *Handler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "missing title", http.StatusBadRequest)
		return
	}

	// Ensure directories exist
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		http.Error(w, "failed to create upload dir", http.StatusInternalServerError)
		return
	}

	thumbDir := filepath.Join(UploadDir, "thumbnails")
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		http.Error(w, "failed to create thumbnail dir", http.StatusInternalServerError)
		return
	}

	id := uuid.New().String()

	rawFilename := id + "_" + header.Filename
	rawPath := filepath.Join(UploadDir, rawFilename)

	finalFilename := id + ".mp4"
	finalPath := filepath.Join(UploadDir, finalFilename)

	thumbName := id + ".jpg"
	thumbPath := filepath.Join(thumbDir, thumbName)

	// Save raw upload
	out, err := os.Create(rawPath)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "failed to write file", http.StatusInternalServerError)
		return
	}

	// Transcode to mp4
	transcodeCmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", rawPath,
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
		http.Error(w, "video processing failed", http.StatusInternalServerError)
		return
	}

	_ = os.Remove(rawPath)

	// Generate thumbnail
	thumbCmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", finalPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-f", "image2",
		thumbPath,
	)

	thumbCmd.Stdout = os.Stdout
	thumbCmd.Stderr = os.Stderr

	if err := thumbCmd.Run(); err != nil {
		http.Error(w, "thumbnail generation failed", http.StatusInternalServerError)
		return
	}

	// Store metadata
	video, err := h.Repo.Create(
		r.Context(),
		id,
		title,
		finalFilename,
		thumbName,
	)
	if err != nil {
		http.Error(w, "database insert failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(video)
}

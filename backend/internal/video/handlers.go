package video

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// Shared upload directory
const UploadDir = "./uploads"

// Handler required by main.go
type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{Repo: repo}
}

// =========================
// GET /videos
// =========================
func (h *Handler) ListVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := h.Repo.List(r.Context())
	if err != nil {
		http.Error(w, "failed to list videos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// =========================
// POST /videos/upload (ASYNC)
// =========================
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
	if err != nil {
	log.Printf("upload error: %v", err)
	http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
	return
}
	defer file.Close()

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "missing title", http.StatusBadRequest)
		return
	}

	// Ensure raw upload directory exists
	rawDir := filepath.Join(UploadDir, "raw")
	if err := os.MkdirAll(rawDir, 0755); err != nil {
		http.Error(w, "failed to create upload dir", http.StatusInternalServerError)
		return
	}

	videoID := uuid.New()
	jobID := uuid.New()


	// Save raw file
	rawFilename := videoID + "_" + header.Filename
	rawPath := filepath.Join(rawDir, rawFilename)

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

	// 1️⃣ Create "processing" video
	video, err := h.Repo.CreateProcessing(
		r.Context(),
		videoID.String(),
		title,
	)
	if err != nil {
		http.Error(w, "failed to create video", http.StatusInternalServerError)
		return
	}

	// 2️⃣ Enqueue background job
	if err := h.Repo.EnqueueJob(
		r.Context(),
		jobID.String(),
		videoID.String(),
		rawPath,
	); err != nil {
		http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
		return
	}

	// 3️⃣ Respond immediately
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(video)
}

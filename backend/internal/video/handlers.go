package video

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{Repo: repo}
}

// POST /videos
func (h *Handler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	video, err := h.Repo.Create(r.Context(), body.Title)
	if err != nil {
		http.Error(w, "failed to create video", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(video)
}

// GET /videos
func (h *Handler) ListVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := h.Repo.List(r.Context())
	if err != nil {
		http.Error(w, "failed to list videos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

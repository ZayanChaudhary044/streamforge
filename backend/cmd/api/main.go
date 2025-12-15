package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ZayanChaudhary044/streamforge/backend/internal/database"
	"github.com/ZayanChaudhary044/streamforge/backend/internal/video"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	db := database.Open(dsn)
	defer db.Close()

	videoRepo := video.NewRepository(db)
	videoHandler := video.NewHandler(videoRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /media/{filename}", video.ServeMedia)
	mux.HandleFunc("GET /thumbnails/{filename}", video.ServeThumbnail)

	mux.HandleFunc("GET /videos", videoHandler.ListVideos)
	mux.HandleFunc("POST /videos/upload", videoHandler.UploadVideo)

	log.Println("backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}

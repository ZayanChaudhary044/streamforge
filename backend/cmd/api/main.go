package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ZayanChaudhary044/streamforge/backend/internal/database"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	db := database.Open(dsn)
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Println("backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

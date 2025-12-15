package video

import (
	"net/http"
	"os"
	"path/filepath"
)

func ServeMedia(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.URL.Path)
	path := filepath.Join("uploads", filename)

	file, err := os.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
}

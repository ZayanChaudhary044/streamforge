package video

import (
	"net/http"
	"path/filepath"
)

func ServeThumbnail(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")

	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeFile(
		w,
		r,
		filepath.Join(UploadDir, "thumbnails", filename),
	)
}

package video

import "time"

type Video struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Filename  *string    `json:"filename"`   // nullable until ready
	Thumbnail *string    `json:"thumbnail"`  // nullable until ready
	Status    string     `json:"status"`     // processing | ready | failed
	Error     *string    `json:"error"`      // nullable
	CreatedAt time.Time  `json:"created_at"`
}

type Job struct {
	ID       string
	VideoID  string
	RawPath  string
	Attempts int
}

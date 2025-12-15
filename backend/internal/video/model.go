package video

import "time"

type Video struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Filename  string    `json:"filename"`
	Thumbnail string     `json:"thumbnail"`
	CreatedAt time.Time `json:"created_at"`

}

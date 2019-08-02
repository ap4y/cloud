package gallery

import "time"

// Album stores album metadata.
type Album struct {
	Name       string    `json:"name"`
	ModTime    time.Time `json:"updated_at"`
	ItemsCount int       `json:"items_count"`
}

// Image stores image metadata.
type Image struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	ModTime time.Time `json:"updated_at"`
}

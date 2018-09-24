package gallery

import "time"

type Album struct {
	Name       string
	ModTime    time.Time
	ItemsCount int
}

type Image struct {
	Name    string
	Path    string
	ModTime time.Time
}

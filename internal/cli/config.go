package cli

import "github.com/ap4y/cloud/api"

// GalleryConfig defines gallery related configuration variables for CLI.
type GalleryConfig struct {
	Path  string `json:"path"`
	Cache string `json:"cache"`
}

// ShareConfig defines share related configuration variables for CLI.
type ShareConfig struct {
	Path string `json:"path"`
}

// Config defines configuration variables for CLI.
type Config struct {
	JWTSecret string            `json:"jwt_secret"`
	Modules   []api.Module      `json:"modules"`
	Users     map[string]string `json:"users"`
	Share     *ShareConfig      `json:"share"`
	Gallery   *GalleryConfig    `json:"gallery"`
}

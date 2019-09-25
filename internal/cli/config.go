package cli

import "github.com/ap4y/cloud/common"

// GalleryConfig defines gallery related configuration variables for CLI.
type GalleryConfig struct {
	Path  string `json:"path"`
	Cache string `json:"cache"`
}

// FilesConfig defines files module related configuration variables for CLI.
type FilesConfig struct {
	Path string `json:"path"`
}

// ShareConfig defines share related configuration variables for CLI.
type ShareConfig struct {
	Path string `json:"path"`
}

// Config defines configuration variables for CLI.
type Config struct {
	JWTSecret string              `json:"jwt_secret"`
	Modules   []common.ModuleType `json:"modules"`
	Users     map[string]string   `json:"users"`
	Share     *ShareConfig        `json:"share"`
	Gallery   *GalleryConfig      `json:"gallery"`
	Files     *FilesConfig        `json:"files"`
}

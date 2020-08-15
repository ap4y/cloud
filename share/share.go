package share

import (
	"github.com/ap4y/cloud/module"
	"github.com/ap4y/cloud/niltime"
)

// Share stores share data.
type Share struct {
	Slug      string       `json:"slug"`
	Type      module.Type  `json:"type"`
	Name      string       `json:"name"`
	Items     []string     `json:"items"`
	ExpiresAt niltime.Time `json:"expires_at"`
}

// IsValid returns true if share is valid.
func (s Share) IsValid() bool {
	if s.Slug == "" {
		return false
	}

	if s.Name == "" {
		return false
	}

	if s.Items == nil || len(s.Items) == 0 {
		return false
	}

	return true
}

// Includes returns true if share includes provided item.
func (s Share) Includes(name, item string) bool {
	if s.Name != name {
		return false
	}

	for _, i := range s.Items {
		if item == i {
			return true
		}
	}

	return false
}

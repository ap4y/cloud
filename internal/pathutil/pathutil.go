package pathutil

import (
	"path/filepath"
	"strings"
)

// Join returns combined path for provided elements. Each element will
// be cleaned before joining by doing:
// - Removing all ".."
func Join(elem ...string) string {
	cleaned := make([]string, len(elem))
	for idx, el := range elem {
		cleaned[idx] = strings.ReplaceAll(el, "..", "")
	}

	return filepath.Join(cleaned...)
}

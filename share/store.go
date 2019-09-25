package share

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Store manages share metadata.
type Store interface {
	// All returns all stores shares.
	All() ([]Share, error)
	// Save persists share metadata.
	Save(share *Share) error
	// Get return share metadata.
	Get(slug string) (*Share, error)
	// Remove removes share metadata.
	Remove(slug string) error
	// Expire removes all expired shares.
	Expire() error
}

type diskStore struct {
	dir string
}

// NewDiskStore returns a new on-disk implementation of the ShareStore.
func NewDiskStore(dir string) (Store, error) {
	if dir == "" {
		return nil, errors.New("dir can't be empty")
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create share dir: %s", err)
		}
	}

	return &diskStore{dir}, nil
}

func (store *diskStore) All() ([]Share, error) {
	path := filepath.Join(store.dir, "*")
	matches, err := filepath.Glob(path)
	if err != nil {
		return nil, fmt.Errorf("file: %s", err)
	}

	shares := make([]Share, 0, len(matches))
	for _, match := range matches {
		_, slug := filepath.Split(match)
		share, err := store.Get(slug)
		if err != nil {
			return nil, err
		}

		shares = append(shares, *share)
	}

	return shares, nil
}

func (store *diskStore) Save(share *Share) error {
	path := filepath.Join(store.dir, share.Slug)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("file: %s", err)
	}

	if err := json.NewEncoder(file).Encode(share); err != nil {
		return fmt.Errorf("json: %s", err)
	}

	return nil
}

func (store *diskStore) Get(slug string) (*Share, error) {
	path := filepath.Join(store.dir, slug)
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("file: %s", err)
	}

	share := &Share{}
	if err := json.NewDecoder(file).Decode(share); err != nil {
		return nil, fmt.Errorf("json: %s", err)
	}

	return share, nil
}

func (store *diskStore) Remove(slug string) error {
	path := filepath.Join(store.dir, slug)
	return os.Remove(path)
}

func (store *diskStore) Expire() error {
	shares, err := store.All()
	if err != nil {
		return err
	}

	for _, share := range shares {
		if share.ExpiresAt.IsZero() {
			continue
		}

		if share.ExpiresAt.After(time.Now()) {
			continue
		}

		if err := store.Remove(share.Slug); err != nil {
			return err
		}
	}

	return nil
}

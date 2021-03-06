package share

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ap4y/cloud/module"
)

func TestShare(t *testing.T) {
	share := &Share{Slug: "bar", Type: module.Gallery, Name: "foo", Items: []string{"test.jpg"}}

	t.Run("IsValid", func(t *testing.T) {
		assert.False(t, Share{}.IsValid())
		assert.False(t, Share{Slug: "bar"}.IsValid())
		assert.False(t, Share{Slug: "bar", Name: "foo"}.IsValid())
		assert.False(t, Share{Slug: "bar", Name: "foo", Items: []string{}}.IsValid())
		assert.True(t, Share{Slug: "bar", Name: "foo", Items: []string{"test.jpg"}}.IsValid())
	})

	t.Run("Includes", func(t *testing.T) {
		assert.False(t, share.Includes("bar", "test.jpg"))
		assert.False(t, share.Includes("foo", "test2.jpg"))
		assert.True(t, share.Includes("foo", "test.jpg"))
	})
}

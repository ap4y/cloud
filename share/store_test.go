package share

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/ap4y/cloud/module"
	"gitlab.com/ap4y/cloud/niltime"
)

func TestShareStore(t *testing.T) {
	dir, err := ioutil.TempDir("", "shares")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store, err := NewDiskStore(dir)
	require.NoError(t, err)

	share := &Share{
		Slug:      "foo",
		Type:      module.Gallery,
		Items:     []string{"foo", "bar"},
		ExpiresAt: niltime.Time{Time: time.Unix(0, 0)},
	}
	t.Run("Save", func(t *testing.T) {
		require.NoError(t, store.Save(share))
	})

	t.Run("All", func(t *testing.T) {
		res, err := store.All()
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, "foo", res[0].Slug)
		assert.Equal(t, module.Gallery, res[0].Type)
		assert.Equal(t, []string{"foo", "bar"}, res[0].Items)
		assert.Equal(t, int64(0), res[0].ExpiresAt.Unix())
	})

	t.Run("Get", func(t *testing.T) {
		res, err := store.Get("foo")
		require.NoError(t, err)
		assert.Equal(t, "foo", res.Slug)
		assert.Equal(t, module.Gallery, res.Type)
		assert.Equal(t, []string{"foo", "bar"}, res.Items)
		assert.Equal(t, int64(0), res.ExpiresAt.Unix())
	})

	t.Run("Remove", func(t *testing.T) {
		require.NoError(t, store.Remove("foo"))

		res, err := store.Get("foo")
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Expire", func(t *testing.T) {
		require.NoError(t, store.Save(share))
		require.NoError(t, store.Save(&Share{Slug: "bar", ExpiresAt: niltime.Time{Time: time.Time{}}}))
		require.NoError(t, store.Expire())

		res, err := store.Get("foo")
		require.Error(t, err)
		assert.Nil(t, res)

		res, err = store.Get("bar")
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

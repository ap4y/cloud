package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiskSource(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	source, err := NewDiskSource(filepath.Join(pwd, "fixtures"))
	require.NoError(t, err)

	t.Run("List", func(t *testing.T) {
		items, err := source.List("")
		require.NoError(t, err)
		require.Len(t, items, 3)

		assert.Equal(t, "foo", items[0].Name)
		assert.Equal(t, ItemTypeFile, items[0].Type)
		assert.Equal(t, "test1", items[1].Name)
		assert.Equal(t, ItemTypeDirectory, items[1].Type)
		assert.Equal(t, "test2", items[2].Name)
		assert.Equal(t, ItemTypeDirectory, items[2].Type)

		items, err = source.List("..")
		require.NoError(t, err)
		require.Len(t, items, 3)
		assert.Equal(t, "foo", items[0].Name)

		items, err = source.List("test1")
		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, "bar", items[0].Name)
		assert.Equal(t, "/test1/bar", items[0].Path)

		items, err = source.List("../../test1")
		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, "bar", items[0].Name)
	})

	t.Run("Content", func(t *testing.T) {
		r, err := source.Content("foo")
		require.NoError(t, err)
		res, err := ioutil.ReadAll(r)
		require.NoError(t, err)
		assert.Equal(t, "foo\n", string(res))

		_, err = source.Content("../foo")
		require.NoError(t, err)
	})

	t.Run("Save", func(t *testing.T) {
		require.NoError(t, source.Save(strings.NewReader("test"), "/test1", "test"))
		defer os.Remove("./fixtures/test1/test")

		res, err := ioutil.ReadFile("./fixtures/test1/test")
		require.NoError(t, err)
		assert.Equal(t, "test", string(res))
	})

	t.Run("Save/unsafe_filename", func(t *testing.T) {
		require.NoError(t, source.Save(strings.NewReader("test"), "/test1", "../test"))
		defer os.Remove("./fixtures/test1/test")

		res, err := ioutil.ReadFile("./fixtures/test1/test")
		require.NoError(t, err)
		assert.Equal(t, "test", string(res))
	})

	t.Run("Save/unsafe_path", func(t *testing.T) {
		require.NoError(t, source.Save(strings.NewReader("test"), "../test1", "../test"))
		defer os.Remove("./fixtures/test1/test")

		res, err := ioutil.ReadFile("./fixtures/test1/test")
		require.NoError(t, err)
		assert.Equal(t, "test", string(res))
	})

	t.Run("Remove", func(t *testing.T) {
		require.NoError(t, source.Save(strings.NewReader("test"), "/test1", "test"))
		require.NoError(t, source.Remove("test1/test"))

		require.NoError(t, source.Save(strings.NewReader("test"), "/test1", "test"))
		require.NoError(t, source.Remove("../test1/test"))
	})
}

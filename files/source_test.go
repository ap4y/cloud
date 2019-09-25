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

	t.Run("Tree", func(t *testing.T) {
		tree, err := source.Tree()
		require.NoError(t, err)

		assert.Equal(t, "/", tree.Name)
		assert.Equal(t, "/", tree.Path)
		assert.Equal(t, ItemTypeDirectory, tree.Type)
		require.Len(t, tree.Children, 3)

		item := tree.Children[0]
		assert.Equal(t, "foo", item.Name)
		assert.Equal(t, "/foo", item.Path)
		assert.Equal(t, ItemTypeFile, item.Type)

		item = tree.Children[1]
		assert.Equal(t, "test1", item.Name)
		assert.Equal(t, "/test1", item.Path)
		assert.Equal(t, ItemTypeDirectory, item.Type)
		require.Len(t, item.Children, 2)
		assert.Equal(t, "bar", item.Children[0].Name)
		assert.Equal(t, "/test1/bar", item.Children[0].Path)
		assert.Equal(t, ItemTypeFile, item.Children[0].Type)

		item = tree.Children[2]
		assert.Equal(t, "test2", item.Name)
		assert.Equal(t, "/test2", item.Path)
		assert.Equal(t, ItemTypeDirectory, item.Type)
		require.Len(t, item.Children, 1)
		assert.Equal(t, "baz", item.Children[0].Name)
		assert.Equal(t, "/test2/baz", item.Children[0].Path)
		assert.Equal(t, ItemTypeFile, item.Children[0].Type)
	})

	t.Run("Mkdir", func(t *testing.T) {
		item, err := source.Mkdir("/test2/folder")
		require.NoError(t, err)
		defer os.RemoveAll("./fixtures/test2/folder")

		assert.Equal(t, ItemTypeDirectory, item.Type)
		assert.Equal(t, "folder", item.Name)
		assert.Equal(t, "/test2/folder", item.Path)
	})

	t.Run("Mkdir/unsafe_path", func(t *testing.T) {
		item, err := source.Mkdir("../test2/folder")
		require.NoError(t, err)
		defer os.RemoveAll("./fixtures/test2/folder")

		assert.Equal(t, ItemTypeDirectory, item.Type)
		assert.Equal(t, "folder", item.Name)
		assert.Equal(t, "/test2/folder", item.Path)
	})

	t.Run("Mkdir/unsafe_folder", func(t *testing.T) {
		item, err := source.Mkdir("/test2/..folder")
		require.NoError(t, err)
		defer os.RemoveAll("./fixtures/test2/folder")

		assert.Equal(t, ItemTypeDirectory, item.Type)
		assert.Equal(t, "folder", item.Name)
		assert.Equal(t, "/test2/folder", item.Path)
	})

	t.Run("File", func(t *testing.T) {
		r, err := source.File("foo")
		require.NoError(t, err)
		res, err := ioutil.ReadAll(r)
		require.NoError(t, err)
		assert.Equal(t, "foo\n", string(res))

		_, err = source.File("../foo")
		require.NoError(t, err)
	})

	t.Run("Save", func(t *testing.T) {
		item, err := source.Save(strings.NewReader("test"), "/test1/test")
		require.NoError(t, err)
		defer os.Remove("./fixtures/test1/test")
		assert.Equal(t, "test", item.Name)
		assert.Equal(t, "/test1/test", item.Path)
		assert.Equal(t, ItemTypeFile, item.Type)

		res, err := ioutil.ReadFile("./fixtures/test1/test")
		require.NoError(t, err)
		assert.Equal(t, "test", string(res))
	})

	t.Run("Save/unsafe_filename", func(t *testing.T) {
		item, err := source.Save(strings.NewReader("test"), "/test1/../test")
		require.NoError(t, err)
		defer os.Remove("./fixtures/test1/test")
		assert.Equal(t, "test", item.Name)
		assert.Equal(t, "/test1/test", item.Path)
		assert.Equal(t, ItemTypeFile, item.Type)

		res, err := ioutil.ReadFile("./fixtures/test1/test")
		require.NoError(t, err)
		assert.Equal(t, "test", string(res))
	})

	t.Run("Save/unsafe_path", func(t *testing.T) {
		item, err := source.Save(strings.NewReader("test"), "../test1/test")
		require.NoError(t, err)
		defer os.Remove("./fixtures/test1/test")
		assert.Equal(t, "test", item.Name)
		assert.Equal(t, "/test1/test", item.Path)
		assert.Equal(t, ItemTypeFile, item.Type)

		res, err := ioutil.ReadFile("./fixtures/test1/test")
		require.NoError(t, err)
		assert.Equal(t, "test", string(res))
	})

	t.Run("Remove", func(t *testing.T) {
		_, err := source.Save(strings.NewReader("test"), "/test1/test")
		require.NoError(t, err)
		item, err := source.Remove("/test1/test")
		require.NoError(t, err)
		assert.Equal(t, "test", item.Name)
		assert.Equal(t, "/test1/test", item.Path)
		assert.Equal(t, ItemTypeFile, item.Type)

		_, err = source.Save(strings.NewReader("test"), "/test1/test")
		require.NoError(t, err)
		item, err = source.Remove("../test1/test")
		require.NoError(t, err)
		assert.Equal(t, "test", item.Name)
		assert.Equal(t, "/test1/test", item.Path)
		assert.Equal(t, ItemTypeFile, item.Type)
	})
}

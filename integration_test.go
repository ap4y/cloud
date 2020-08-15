package integration

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ap4y/cloud/api"
	"github.com/ap4y/cloud/files"
	"github.com/ap4y/cloud/gallery"
	"github.com/ap4y/cloud/module"
	"github.com/ap4y/cloud/share"
)

var privateRoutes = []struct {
	method string
	url    string
	body   string
	form   bool
}{
	{"GET", "/modules", "", false},
	{"GET", "/shares", "", false},
	{"POST", "/shares", "{\"type\":\"gallery\",\"name\":\"foo\",\"items\":[\"test.jpg\"]}", false},
	{"DELETE", "/shares/foo", "", false},
	{"GET", "/gallery", "", false},
	{"GET", "/gallery/album1/images", "", false},
	{"GET", "/gallery/album1/image/test.jpg", "", false},
	{"GET", "/gallery/album1/thumbnail/test.jpg", "", false},
	{"GET", "/gallery/album1/exif/test.jpg", "", false},
	{"GET", "/files", "", false},
	{"POST", "/files/mkdir/testfoo", "", false},
	{"POST", "/files/rmdir/testfoo", "", false},
	{"POST", "/files/upload/test1", "bar", true},
	{"DELETE", "/files/file/test1/test", "", false},
	{"GET", "/files/file/foo", "", false},
}

var publicRoutes = []struct {
	method string
	url    string
	body   string
}{
	{"POST", "/user/sign_in", "{\"username\":\"test\",\"password\":\"changeme\"}"},
	{"GET", "/share/bar", ""},
	{"GET", "/share/bar/gallery/album1/images", ""},
	{"GET", "/share/bar/gallery/album1/image/test.jpg", ""},
	{"GET", "/share/bar/gallery/album1/thumbnail/test.jpg", ""},
	{"GET", "/share/bar/gallery/album1/exif/test.jpg", ""},
	{"GET", "/share/baz/files", ""},
	{"GET", "/share/baz/files/file/test1/inner/foo", ""},
}

var prohibitedRoutes = []struct {
	method string
	url    string
	body   string
}{
	{"GET", "/share/bar/gallery", ""},
	{"POST", "/share/baz/files/mkdir/testfoo", ""},
	{"POST", "/share/baz/files/rmdir/testfoo", ""},
	{"POST", "/share/baz/files/upload/foo", ""},
	{"DELETE", "/share/baz/files/file/foo", ""},
}

func TestAPIServer(t *testing.T) {
	cacheDir, err := ioutil.TempDir("", "cache")
	require.NoError(t, err)
	defer os.RemoveAll(cacheDir)

	modules := map[module.Type]http.Handler{
		module.Gallery: galleryModule(t, cacheDir),
		module.Files:   filesModule(t),
	}
	cs := api.NewMemoryCredentialsStorage(
		map[string]string{"test": "$2b$10$fEWhY87kzeaV3hUEB6phTuyWjpv73V5m.YcqTxHXnvqEGIou1tXGO"},
		jwt.SigningMethodHS256,
		[]byte("secret"),
	)

	sharesDir, err := ioutil.TempDir("", "shares")
	require.NoError(t, err)
	defer os.RemoveAll(sharesDir)
	ss, err := share.NewDiskStore(sharesDir)
	require.NoError(t, err)
	err = ss.Save(&share.Share{Slug: "foo", Type: module.Gallery, Name: "foo", Items: []string{"test.jpg"}})
	require.NoError(t, err)
	err = ss.Save(&share.Share{Slug: "bar", Type: module.Gallery, Name: "album1", Items: []string{"test.jpg"}})
	require.NoError(t, err)
	err = ss.Save(&share.Share{Slug: "baz", Type: module.Files, Name: "/test1", Items: []string{"/test1/inner"}})
	require.NoError(t, err)

	handler, err := api.NewServer(modules, cs, ss)
	require.NoError(t, err)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	client := ts.Client()

	for _, tc := range privateRoutes {
		t.Run(fmt.Sprintf("anonymous/%s%s", tc.method, tc.url), func(t *testing.T) {
			var body io.Reader
			contentType := "application/json"
			if tc.form {
				body, contentType = formFile(t, "test", tc.body)
			} else {
				body = strings.NewReader(tc.body)
			}
			req, err := http.NewRequest(tc.method, ts.URL+"/api"+tc.url, body)
			require.NoError(t, err)
			req.Header.Set("Content-Type", contentType)
			res, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		})
	}

	jwtToken := jwtToken("test", "secret")
	for _, tc := range privateRoutes {
		t.Run(fmt.Sprintf("%s%s", tc.method, tc.url), func(t *testing.T) {
			var body io.Reader
			contentType := "application/json"
			if tc.form {
				body, contentType = formFile(t, "test", tc.body)
			} else {
				body = strings.NewReader(tc.body)
			}
			req, err := http.NewRequest(tc.method, ts.URL+"/api"+tc.url, body)
			require.NoError(t, err)
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Cookie", "token="+jwtToken+";")
			res, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
		})
	}

	for _, tc := range publicRoutes {
		t.Run(fmt.Sprintf("%s%s", tc.method, tc.url), func(t *testing.T) {
			req, err := http.NewRequest(tc.method, ts.URL+"/api"+tc.url, strings.NewReader(tc.body))
			require.NoError(t, err)
			res, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
		})
	}

	for _, tc := range prohibitedRoutes {
		t.Run(fmt.Sprintf("%s%s", tc.method, tc.url), func(t *testing.T) {
			req, err := http.NewRequest(tc.method, ts.URL+"/api"+tc.url, strings.NewReader(tc.body))
			require.NoError(t, err)
			res, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, res.StatusCode)
		})
	}
}

func galleryModule(t *testing.T, cacheDir string) http.Handler {
	t.Helper()

	pwd, err := os.Getwd()
	require.NoError(t, err)
	source, err := gallery.NewDiskSource(filepath.Join(pwd, "/gallery/fixtures"), []string{".jpg"})
	require.NoError(t, err)

	cache, err := gallery.NewDiskCache(cacheDir)
	require.NoError(t, err)

	return gallery.NewGalleryAPI(source, cache)
}

func filesModule(t *testing.T) http.Handler {
	t.Helper()

	pwd, err := os.Getwd()
	require.NoError(t, err)
	source, err := files.NewDiskSource(filepath.Join(pwd, "/files/fixtures"))
	require.NoError(t, err)

	return files.NewFilesAPI(source)
}

func jwtToken(username, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{api.UserAuthKey: username})
	str, _ := token.SignedString([]byte(secret))
	return str
}

func formFile(t *testing.T, name, content string) (io.Reader, string) {
	t.Helper()

	buf := new(bytes.Buffer)
	formWriter := multipart.NewWriter(buf)
	fw, err := formWriter.CreateFormFile("file", name)
	require.NoError(t, err)
	_, err = fw.Write([]byte(content))
	require.NoError(t, err)
	formWriter.Close()
	return buf, formWriter.FormDataContentType()
}

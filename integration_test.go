package integration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ap4y/cloud/api"
	"github.com/ap4y/cloud/common"
	"github.com/ap4y/cloud/gallery"
	"github.com/ap4y/cloud/share"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var privateRoutes = []struct {
	method string
	url    string
	body   string
}{
	{"GET", "/modules", ""},
	{"GET", "/shares", ""},
	{"POST", "/shares", "{\"type\":\"gallery\",\"name\":\"foo\",\"items\":[\"test.jpg\"]}"},
	{"DELETE", "/shares/foo", ""},
	{"GET", "/gallery", ""},
	{"GET", "/gallery/album1/images", ""},
	{"GET", "/gallery/album1/image/test.jpg", ""},
	{"GET", "/gallery/album1/thumbnail/test.jpg", ""},
	{"GET", "/gallery/album1/exif/test.jpg", ""},
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
}

func TestAPIServer(t *testing.T) {
	cacheDir, err := ioutil.TempDir("", "cache")
	require.NoError(t, err)
	defer os.RemoveAll(cacheDir)

	modules := map[common.ModuleType]http.Handler{common.ModuleGallery: galleryModule(t, cacheDir)}
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
	err = ss.Save(&share.Share{Slug: "foo", Type: common.ModuleGallery, Name: "foo", Items: []string{"test.jpg"}})
	require.NoError(t, err)
	err = ss.Save(&share.Share{Slug: "bar", Type: common.ModuleGallery, Name: "album1", Items: []string{"test.jpg"}})
	require.NoError(t, err)

	handler, err := api.NewServer(modules, cs, ss)
	require.NoError(t, err)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	client := ts.Client()

	for _, tc := range privateRoutes {
		t.Run(fmt.Sprintf("anonymous/%s%s", tc.method, tc.url), func(t *testing.T) {
			req, err := http.NewRequest(tc.method, ts.URL+"/api"+tc.url, strings.NewReader(tc.body))
			require.NoError(t, err)
			res, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		})
	}

	jwtToken := jwtToken("test", "secret")
	for _, tc := range privateRoutes {
		t.Run(fmt.Sprintf("%s%s", tc.method, tc.url), func(t *testing.T) {
			req, err := http.NewRequest(tc.method, ts.URL+"/api"+tc.url, strings.NewReader(tc.body))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+jwtToken)
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

func jwtToken(username, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{api.UserAuthKey: username})
	str, _ := token.SignedString([]byte(secret))
	return str
}

package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/ap4y/cloud/contextkey"
)

func TestAuth(t *testing.T) {
	credentials := NewMemoryCredentialsStorage(
		map[string]string{"test": "$2b$10$fEWhY87kzeaV3hUEB6phTuyWjpv73V5m.YcqTxHXnvqEGIou1tXGO"},
		jwt.SigningMethodHS256,
		[]byte("secret"),
	)

	jwtToken := func(username string) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{UserAuthKey: username})
		str, _ := token.SignedString([]byte("secret"))
		return str
	}

	t.Run("AuthHandler", func(t *testing.T) {
		tcs := []struct {
			name     string
			username string
			password string
			status   int
			res      string
		}{
			{"unknown user", "foo", "bar", http.StatusBadRequest, "{\"error\":\"Failed to authenticate user: invalid username or password\"}\n"},
			{"invalid password", "test", "bar", http.StatusBadRequest, "{\"error\":\"Failed to authenticate user: invalid username or password\"}\n"},
			{"valid", "test", "changeme", http.StatusOK, "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoidGVzdCJ9.wk1swko8GbuwRMRuTR6q_x7AZQGbbwm8sZLyg90afbs\"}\n"},
		}

		api := AuthHandler(credentials)
		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				body := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\"}", tc.username, tc.password)
				req := httptest.NewRequest("POST", "http://cloud.api/sign_in", strings.NewReader(body))

				api.ServeHTTP(w, req)
				resp := w.Result()
				require.Equal(t, tc.status, resp.StatusCode)

				res, _ := ioutil.ReadAll(resp.Body)
				assert.Equal(t, tc.res, string(res))
			})
		}
	})

	t.Run("Authenticator", func(t *testing.T) {
		tcs := []struct {
			name     string
			token    string
			status   int
			username string
		}{
			{"invalid token", "foo", http.StatusUnauthorized, ""},
			{"invalid username", jwtToken("foo"), http.StatusUnauthorized, ""},
			{"valid", jwtToken("test"), http.StatusOK, "test"},
		}

		for _, tc := range tcs {
			t.Run("header - "+tc.name, func(t *testing.T) {
				handler := Authenticator(credentials)(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						assert.Equal(t, tc.username, r.Context().Value(contextkey.UsernameCtxKey))
						io.WriteString(w, "<html><body>Hello World!</body></html>") // nolint: errcheck
					}),
				)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "http://cloud.api", nil)
				req.Header.Set("Authorization", "Bearer "+tc.token)
				handler.ServeHTTP(w, req)

				resp := w.Result()
				require.Equal(t, tc.status, resp.StatusCode)
			})

			t.Run("query - "+tc.name, func(t *testing.T) {
				handler := Authenticator(credentials)(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						assert.Equal(t, tc.username, r.Context().Value(contextkey.UsernameCtxKey))
						io.WriteString(w, "<html><body>Hello World!</body></html>") // nolint: errcheck
					}),
				)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "http://cloud.api?jwt="+tc.token, nil)
				handler.ServeHTTP(w, req)

				resp := w.Result()
				require.Equal(t, tc.status, resp.StatusCode)
			})
		}
	})
}

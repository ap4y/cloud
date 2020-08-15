package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"

	"github.com/ap4y/cloud/contextkey"
	"github.com/ap4y/cloud/internal/httputil"
)

// UserAuthKey defines usename key in jwt token.
const UserAuthKey = "user"

const tokenCookieKey = "token"

// CredentialsStorage stores and validates user credentials.
type CredentialsStorage interface {
	// Authenticate returns jwt token if provided password matches to a
	// stored hash for a given user, error returned otherwise.
	Authenticate(username, password string) (tokenString string, err error)
	// Validate validates provided jwt token against stored credentials.
	Validate(tokenString string) (username string, err error)
}

type memoryCredentialsStorage struct {
	hashes        map[string]string
	signingMethod jwt.SigningMethod
	signKey       interface{}
}

// NewMemoryCredentialsStorage returns a new CredentialsStorage that stores user credentials in memory.
func NewMemoryCredentialsStorage(hashes map[string]string, signingMethod jwt.SigningMethod, signKey interface{}) CredentialsStorage {
	return &memoryCredentialsStorage{hashes, signingMethod, signKey}
}

func (cs *memoryCredentialsStorage) Authenticate(username, password string) (string, error) {
	hashedPassword := cs.hashes[username]
	if hashedPassword == "" {
		return "", fmt.Errorf("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	token := jwt.NewWithClaims(cs.signingMethod, jwt.MapClaims{UserAuthKey: username})
	tokenString, err := token.SignedString(cs.signKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %s", err)
	}

	return tokenString, nil
}

func (cs *memoryCredentialsStorage) Validate(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != cs.signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return cs.signKey, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	username, ok := claims[UserAuthKey].(string)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	if hashedPassword := cs.hashes[username]; hashedPassword == "" {
		return "", fmt.Errorf("invalid token claims")
	}

	return username, nil
}

// AuthHandler returns a new handler for authentication endpoints.
func AuthHandler(credentials CredentialsStorage) http.Handler {
	mux := chi.NewRouter()
	mux.Post("/sign_in", func(w http.ResponseWriter, req *http.Request) {
		body := map[string]string{}

		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			httputil.Error(w, fmt.Sprintf("Failed to decode json: %s", err), http.StatusBadRequest)
			return
		}

		token, err := credentials.Authenticate(body["username"], body["password"])
		if err != nil {
			httputil.Error(w, fmt.Sprintf("Failed to authenticate user: %s", err), http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     tokenCookieKey,
			Value:    token,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
		})
		httputil.Respond(w, map[string]string{"token": token})
	})

	return mux
}

// Authenticator returns authentication middleware.
func Authenticator(credentials CredentialsStorage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token, err := req.Cookie(tokenCookieKey)
			if err != nil {
				httputil.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			username, err := credentials.Validate(token.Value)
			if err != nil {
				httputil.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(req.Context(), contextkey.UsernameCtxKey, username)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

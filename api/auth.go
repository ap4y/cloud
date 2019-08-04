package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

// UserAuthKey defines usename key in jwt token.
const UserAuthKey = "user"

// UsernameCtxKey defines username request context key.
var UsernameCtxKey = &contextKey{"Username"}

type contextKey struct {
	name string
}

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
	mux.Post("/sign_in", func(w http.ResponseWriter, res *http.Request) {
		body := map[string]string{}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode json: %s", err), http.StatusBadRequest)
			return
		}

		token, err := credentials.Authenticate(body["username"], body["password"])
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to authenticate user: %s", err), http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode json: %s", err), http.StatusBadRequest)
		}
	})

	return mux
}

// Authenticator returns authentication middleware.
func Authenticator(credentials CredentialsStorage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := tokenFromHeader(r.Header.Get("Authorization"))
			if token == "" {
				token = r.URL.Query().Get("jwt")
			}

			username, err := credentials.Validate(token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UsernameCtxKey, username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func tokenFromHeader(header string) string {
	if header == "" {
		return ""
	}

	if !strings.HasPrefix(strings.ToLower(header), "bearer") {
		return ""
	}

	return header[7:]
}

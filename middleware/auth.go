package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")

		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, errors.New("auth error: no token").Error(), http.StatusUnauthorized)
				return
			}

			token := cookie.Value
			jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				secret := os.Getenv("TODO_TOKEN_SECRET")
				if len(secret) == 0 {
					return nil, errors.New("sekret key not found")
				}
				return []byte(secret), nil
			})
			if err != nil {
				http.Error(w, "token parsing error: "+err.Error(), http.StatusUnauthorized)
				return
			}
			claims, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, errors.New("token parsing error").Error(), http.StatusUnauthorized)
				return
			}
			passwordHash, ok := claims["passwordHash"].(string)
			if !ok {
				http.Error(w, errors.New("token parsing error").Error(), http.StatusUnauthorized)
				return
			}
			h := sha256.Sum256([]byte(pass))
			hString := hex.EncodeToString(h[:])
			if passwordHash != hString {
				http.Error(w, errors.New("token validation error").Error(), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

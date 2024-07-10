package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/IvanMishnev/go_final_project/internal/constants"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := constants.Password

		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, errors.New("auth error: no token").Error(), http.StatusUnauthorized)
				return
			}

			token := cookie.Value
			err = validateToken(token, pass)
			if err != nil {
				http.Error(w, "token validation error: "+err.Error(), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func validateToken(token string, pass string) error {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		secret := constants.TockenSecret
		if len(secret) == 0 {
			return nil, errors.New("sekret key not found")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return err
	}

	passwordHash, ok := claims["passwordHash"].(string)
	if !ok {
		return err
	}

	h := sha256.Sum256([]byte(pass))
	hString := hex.EncodeToString(h[:])
	if passwordHash != hString {
		return err
	}

	return nil
}

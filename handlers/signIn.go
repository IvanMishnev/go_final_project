package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	mPas := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&mPas)
	if err != nil {
		JSONError(w, "JSON deserealization error: "+err.Error(), http.StatusBadRequest)
		return
	}
	password, ok := mPas["password"]
	if !ok {
		JSONError(w, errors.New("wrong JSON format").Error(), http.StatusBadRequest)
		return
	}

	if password == os.Getenv("TODO_PASSWORD") {
		secret := os.Getenv("TODO_TOKEN_SECRET")
		pHash := sha256.Sum256([]byte(password))
		passwordHashString := hex.EncodeToString(pHash[:])
		claims := jwt.MapClaims{
			"passwordHash": passwordHashString,
		}
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := jwtToken.SignedString([]byte(secret))
		if err != nil {
			JSONError(w, "token creating error: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token": signedToken,
		})
	} else {
		JSONError(w, errors.New("wrong password").Error(), http.StatusBadRequest)
		return
	}
}

func FileServer(w http.ResponseWriter, r *http.Request) {
	dir := "./web"
	http.ServeFile(w, r, dir+r.URL.Path)
}

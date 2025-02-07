package handlers

import (
	"encoding/json"
	"net/http"
)

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	mErr := map[string]interface{}{
		"error": err,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(mErr)
}

package controller

import (
	"encoding/json"
	"net/http"
	"product-api/internal/jwt"
)

func GenerateJWTToken(w http.ResponseWriter, r *http.Request) {

	// Generate the JWT token
	token, err := jwt.GenerateToken()
	if err != nil {
		http.Error(w, "Error generating token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the generated token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

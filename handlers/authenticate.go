package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/o1egl/paseto"
)

// AuthRequest represents a request to authenticate a user.
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents a response containing an authentication token.
type AuthResponse struct {
	Token string `json:"token"`
}

// AuthenticateHandler handles authentication requests.
func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Implement actual authentication logic here.
	if req.Username != "user" || req.Password != "password" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate a PASETO token with the user's username as a claim.
	payload := paseto.JSONToken{
		Subject: req.Username,
	}
	token, err := paseto.Encrypt(encryptionKey, payload, signingKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := AuthResponse{
		Token: token,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
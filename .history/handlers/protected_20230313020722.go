package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/o1egl/paseto"
)

// ProtectedHandler handles requests to a protected endpoint.
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// Get the token from the Authorization header.
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}
	token, err := paseto.ParseToken([]byte(authHeader), encryptionKey, signingKey)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Get the username claim from the token.
	var payload paseto.JSONToken
	if err := token.GetPayload(&payload); err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, " + payload.Subject + "!",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AuthenticationMiddleware is a middleware that checks for a valid token in the Authorization header.
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		token, err := paseto.ParseToken([]byte(authHeader), encryptionKey, signingKey)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Get the username claim from the token.
		var payload paseto.JSONToken
		if err := token.GetPayload(&payload); err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add the username claim to the request context.
		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", payload.Subject)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
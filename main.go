package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yourusername/yourprojectname/handlers"
)

func main() {
	r := mux.NewRouter()

	// Public endpoints
	r.HandleFunc("/authenticate", handlers.AuthenticateHandler).Methods("POST")

	// Protected endpoints
	protected := r.PathPrefix("/protected").Subrouter()
	protected.Use(handlers.AuthenticationMiddleware)
	protected.HandleFunc("", handlers.ProtectedHandler).Methods("GET")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

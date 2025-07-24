// File: internal/handler/cors.go
package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func WithCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Printf("Error loading .env file")
		}

		HOST := os.Getenv("HOST")

		// Allow your frontend origin here
		w.Header().Set("Access-Control-Allow-Origin", HOST)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h(w, r)
	}
}

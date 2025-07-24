// File: internal/handler/cors.go
package handler

import (
	"log"
	"net/http"
	"os"
)

func WithCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		hostFromEnv := os.Getenv("HOST")

		// Log all values
		log.Println("🧠 Debug CORS Middleware")
		log.Printf("📥 Request Origin: %q\n", origin)
		log.Printf("📦 HOST from .env: %q\n", hostFromEnv)

		// Check if they match exactly
		if origin == hostFromEnv {
			log.Println("✅ Origin matches HOST from .env — setting CORS header.")
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			log.Println("❌ Origin does NOT match HOST — CORS header NOT set.")
		}

		// Always log what CORS headers were actually set
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Vary", "Origin")

		// Handle OPTIONS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	}
}

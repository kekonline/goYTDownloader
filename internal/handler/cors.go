// File: internal/handler/cors.go
package handler

import (
	"log"
	"net/http"
	"os"
)

func WithCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("🧠 Debug CORS Middleware")

		origin := r.Header.Get("Origin")
		log.Printf("📥 Request Origin: %q\n", origin)

		host := os.Getenv("HOST")
		log.Printf("📦 HOST from .env: %q\n", host)

		if origin == host {
			log.Println("✅ Origin matches HOST from .env — setting CORS headers")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			log.Println("❌ Origin did not match — no CORS headers set")
		}

		// Handle preflight request
		if r.Method == http.MethodOptions {
			log.Println("🚦 OPTIONS Preflight request — returning early")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Continue with normal handler
		h(w, r)
	}
}

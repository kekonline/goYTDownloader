// File: internal/handler/cors.go
package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func WithCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load .env file
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Warning: .env file not loaded: %v\n", err)
		}

		HOST := os.Getenv("HOST")

		// log.Printf("CORS middleware: HOST=%s", HOST)

		// Got from request headers
		// log.Printf("CORS middleware: Origin=%s", r.Header.Get("Origin"))
		// log.Printf("CORS middleware: Method=%s", r.Method)
		// log.Printf("CORS middleware: Headers=%v", r.Header)
		// log.Printf("CORS middleware: URL=%s", r.URL)
		// log.Printf("CORS middleware: RemoteAddr=%s", r.RemoteAddr)
		// log.Printf("CORS middleware: RequestURI=%s", r.RequestURI)
		// log.Printf("CORS middleware: Host=%s", r.Host)

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

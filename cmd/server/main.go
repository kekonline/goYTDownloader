package main

import (
	"log"
	"net/http"
	"os"

	"goYTDownloader/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" // fallback for local dev
	}

	// POST endpoints wrapped with CORS middleware
	http.HandleFunc("/api/audio-stream", handler.WithCORS(handler.AudioStreamHandler))
	http.HandleFunc("/api/audio-stream-multiple-files-poc", handler.WithCORS(handler.AudioStreamHandlerMultipleFilesPoc))
	http.HandleFunc("/api/audio-streamv2", handler.WithCORS(handler.AudioStreamHandlerV2))

	// OPTIONS handlers for preflight requests on the same endpoints
	http.HandleFunc("/api/audio-stream", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handler.WithCORS(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})(w, r)
			return
		}
		handler.WithCORS(handler.AudioStreamHandler)(w, r)
	})
	http.HandleFunc("/api/audio-stream-multiple-files-poc", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handler.WithCORS(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})(w, r)
			return
		}
		handler.WithCORS(handler.AudioStreamHandlerMultipleFilesPoc)(w, r)
	})
	http.HandleFunc("/api/audio-streamv2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handler.WithCORS(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})(w, r)
			return
		}
		handler.WithCORS(handler.AudioStreamHandlerV2)(w, r)
	})

	log.Println("Server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

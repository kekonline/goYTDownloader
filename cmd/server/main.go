package main

import (
	"fmt"
	"goYTDownloader/internal/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from goYTDownloader!")
	})
	http.HandleFunc("/api/audio-stream", handler.AudioStreamHandler)
	http.HandleFunc("/api/audio-stream-with-cors", handler.WithCORS(handler.AudioStreamHandler))
	http.HandleFunc("/api/audio-stream-multiple-files-poc", handler.WithCORS(handler.AudioStreamHandlerMultipleFilesPoc))
	http.HandleFunc("/api/audio-streamv2", handler.WithCORS(handler.AudioStreamHandlerV2))

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT environment variable not set, using default port 10000")
		port = "10000" // Fallback for local dev
	}

	log.Println("Server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

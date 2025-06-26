package main

import (
	"log"
	"net/http"

	"goYTDownloader/internal/handler"
)

// create . env file with server name

func main() {
	http.HandleFunc("/api/audio-stream", handler.WithCORS(handler.AudioStreamHandler))
	http.HandleFunc("/api/audio-stream-multiple-files-poc", handler.WithCORS(handler.AudioStreamHandlerMultipleFilesPoc))
	// http.HandleFunc("/api/audio-stream", handler.AudioStreamHandler)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

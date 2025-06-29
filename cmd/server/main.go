package main

import (
	"log"
	"net/http"

	"goYTDownloader/internal/handler"
)

func main() {
	http.HandleFunc("/api/audio-stream", handler.WithCORS(handler.AudioStreamHandler))
	http.HandleFunc("/api/audio-stream-multiple-files-poc", handler.WithCORS(handler.AudioStreamHandlerMultipleFilesPoc))
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

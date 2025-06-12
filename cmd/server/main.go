package main

import (
	"log"
	"net/http"

	"goYTDownloader/internal/handler"
)

func main() {
	http.HandleFunc("/api/audio-stream", handler.AudioStreamHandler)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

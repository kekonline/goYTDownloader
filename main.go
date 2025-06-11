package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
)

type DownloadRequest struct {
	URL string `json:"url"`
}

func handleAudioStream(w http.ResponseWriter, r *http.Request) {
	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Run yt-dlp command: best audio, output to stdout
	// Adjust format (-f) or args if needed
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-o", "-", req.URL)

	// Get stdout pipe to stream output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Failed to get yt-dlp output", http.StatusInternalServerError)
		return
	}

	// Start command
	if err := cmd.Start(); err != nil {
		http.Error(w, "Failed to start yt-dlp", http.StatusInternalServerError)
		return
	}

	// Set headers for audio stream (change Content-Type if needed)
	w.Header().Set("Content-Type", "audio/webm") // or audio/mpeg, etc.
	w.Header().Set("Content-Disposition", "attachment; filename=\"audio.webm\"")

	// Stream yt-dlp stdout directly to response
	_, err = io.Copy(w, stdout)
	if err != nil {
		log.Printf("Error streaming audio: %v", err)
	}

	// Wait for yt-dlp to finish
	if err := cmd.Wait(); err != nil {
		log.Printf("yt-dlp command error: %v", err)
	}
}

func main() {
	http.HandleFunc("/api/audio-stream", handleAudioStream)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

type DownloadRequest struct {
	URL string `json:"url"`
}

func handleAudioStreamURL(w http.ResponseWriter, r *http.Request) {
	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-g", req.URL)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("yt-dlp error: %s\n", err.Error())
		http.Error(w, "Failed to get stream URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"stream_url": string(output),
	})
}

func main() {
	http.HandleFunc("/api/audio-stream", handleAudioStreamURL)
	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

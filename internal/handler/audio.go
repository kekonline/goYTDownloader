package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"

	"goYTDownloader/internal/model"
)

func AudioStreamHandler(w http.ResponseWriter, r *http.Request) {
	var req model.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-o", "-", req.URL)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Failed to get yt-dlp output", http.StatusInternalServerError)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("GOT TO HERe")
		http.Error(w, "Failed to start yt-dlp", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/webm")
	w.Header().Set("Content-Disposition", "attachment; filename=\"audio.webm\"")

	_, err = io.Copy(w, stdout)
	if err != nil {
		log.Printf("Error streaming audio: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("yt-dlp command error: %v", err)
	}
}

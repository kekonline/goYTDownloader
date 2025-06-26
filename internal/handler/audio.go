package handler

import (
	"encoding/json"
	"fmt"
	"goYTDownloader/internal/model"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
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
		http.Error(w, "Failed to start yt-dlp", http.StatusInternalServerError)
		return
	}

	cmdName := exec.Command("yt-dlp", "--get-filename", "-o", "%(title)s.%(ext)s", req.URL)
	output, err := cmdName.Output()
	if err != nil {
		http.Error(w, "Failed to get filename", http.StatusInternalServerError)
		return
	}
	filename := strings.TrimSpace(string(output))
	log.Printf("Streaming audio for: %s", filename)

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	HOST := os.Getenv("HOST")

	origin := r.Header.Get("Origin")
	if origin == HOST {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	}
	w.Header().Set("Content-Type", "audio/webm")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	_, err = io.Copy(w, stdout)
	if err != nil {
		log.Printf("Error streaming audio: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("yt-dlp command error: %v", err)
	}
}

func AudioStreamHandlerMultipleFilesPoc(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to start yt-dlp", http.StatusInternalServerError)
		return
	}

	cmdName := exec.Command("yt-dlp", "--get-filename", "-o", "%(title)s.%(ext)s", req.URL)
	output, err := cmdName.Output()
	if err != nil {
		http.Error(w, "Failed to get filename", http.StatusInternalServerError)
		return
	}
	filename := strings.TrimSpace(string(output))
	log.Printf("Streaming audio for: %s", filename)

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	HOST := os.Getenv("HOST")

	origin := r.Header.Get("Origin")
	if origin == HOST {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	}
	w.Header().Set("Content-Type", "audio/webm")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	_, err = io.Copy(w, stdout)
	if err != nil {
		log.Printf("Error streaming audio: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("yt-dlp command error: %v", err)
	}
}

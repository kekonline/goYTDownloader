package handler

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"goYTDownloader/internal/model"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

const (
	maxStreamSize = 50 * 1024 * 1024 // 50MB
	timeout       = 60 * time.Second // Timeout for entire operation
)

func AudioStreamHandlerV2(w http.ResponseWriter, r *http.Request) {
	// Apply a timeout to the whole request
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	var req model.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		fmt.Printf("Invalid JSON request: %v\n", err)
		return
	}

	// Check if yt-dlp is installed
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		http.Error(w, "yt-dlp is not installed or not in PATH", http.StatusServiceUnavailable)
		fmt.Printf("yt-dlp not found: %v\n", err)
		return
	}

	// Get the filename
	cmdName := exec.CommandContext(ctx, "yt-dlp", "--get-filename", "-o", "%(title)s.%(ext)s", req.URL)
	output, err := cmdName.Output()
	if err != nil {
		http.Error(w, "Failed to get filename", http.StatusInternalServerError)
		fmt.Printf("Failed to get filename: %v\n", err)
		return
	}
	filename := strings.TrimSpace(string(output))
	fmt.Printf("Streaming audio for: %s\n", filename)

	// Start yt-dlp streaming audio to stdout
	cmd := exec.CommandContext(ctx, "yt-dlp", "-f", "bestaudio", "-o", "-", req.URL)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Failed to get yt-dlp output", http.StatusInternalServerError)
		fmt.Printf("StdoutPipe error: %v\n", err)
		return
	}
	if err := cmd.Start(); err != nil {
		http.Error(w, "Failed to start yt-dlp", http.StatusInternalServerError)
		fmt.Printf("Start error: %v\n", err)
		return
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not loaded: %v\n", err)
	}
	HOST := os.Getenv("HOST")
	if origin := r.Header.Get("Origin"); origin == HOST {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	}

	w.Header().Set("Content-Type", "audio/webm")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// Limit the amount of data streamed
	limitedReader := io.LimitReader(stdout, maxStreamSize)
	_, err = io.Copy(w, limitedReader)
	if err != nil {
		fmt.Printf("Streaming error: %v\n", err)
	}

	// Wait for yt-dlp to finish
	if err := cmd.Wait(); err != nil {
		fmt.Printf("yt-dlp finished with error: %v\n", err)
	}
}

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

	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not loaded: %v", err)
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

type Song struct {
	Title string
	Data  []byte
}

func AudioStreamHandlerMultipleFilesPoc(w http.ResponseWriter, r *http.Request) {
	numWorkers := 3

	var req model.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	urls := strings.Split(req.URL, ",")
	numJobs := len(urls)

	if numJobs == 0 {
		http.Error(w, "No URLs provided", http.StatusBadRequest)
		return
	}

	log.Printf("Received URLs: %v", urls)

	jobs := make(chan string, numJobs)
	results := make(chan Song, numJobs)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go Worker(i, jobs, results, &wg)
	}

	go func() {
		for _, url := range urls {
			jobs <- strings.TrimSpace(url)
		}
		close(jobs)
	}()

	wg.Wait()
	close(results)

	// Load .env and set headers
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file")
	}
	HOST := os.Getenv("HOST")
	origin := r.Header.Get("Origin")
	if origin == HOST {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="audio_files.zip"`)

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	filesWritten := 0
	for result := range results {
		if result.Data == nil {
			log.Printf("Skipping empty result: %s", result.Title)
			continue
		}

		f, err := zipWriter.Create(result.Title)
		if err != nil {
			log.Printf("Failed to create zip entry for %s: %v", result.Title, err)
			continue
		}

		_, err = f.Write(result.Data)
		if err != nil {
			log.Printf("Failed to write data for %s: %v", result.Title, err)
			continue
		}

		filesWritten++
	}

	if filesWritten == 0 {
		http.Error(w, "All downloads failed", http.StatusInternalServerError)
	}
}

func Worker(id int, jobs <-chan string, results chan<- Song, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		log.Printf("Worker %d: downloading %s", id, job)

		// Added --no-playlist to avoid playlist hangs; remove if you want playlist support
		cmd := exec.Command("yt-dlp", "--no-playlist", "-f", "bestaudio", "-o", "-", job)
		cmd.Stdin = nil

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("Worker %d: Failed to get stdout for %s: %v", id, job, err)
			results <- Song{Title: "", Data: nil}
			continue
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Printf("Worker %d: Failed to get stderr for %s: %v", id, job, err)
			results <- Song{Title: "", Data: nil}
			continue
		}

		if err := cmd.Start(); err != nil {
			log.Printf("Worker %d: Failed to start yt-dlp for %s: %v", id, job, err)
			results <- Song{Title: "", Data: nil}
			continue
		}

		var stderrBuf bytes.Buffer
		go func() {
			io.Copy(&stderrBuf, stderr)
		}()

		// Get the filename
		cmdName := exec.Command("yt-dlp", "--no-playlist", "--get-filename", "-o", "%(title)s.%(ext)s", job)
		output, err := cmdName.Output()
		if err != nil {
			log.Printf("Worker %d: Failed to get filename for %s: %v", id, job, err)
			results <- Song{Title: "", Data: nil}
			_ = cmd.Wait()
			continue
		}
		filename := strings.TrimSpace(string(output))

		// Read audio stream
		data, err := io.ReadAll(stdout)
		if err != nil {
			log.Printf("Worker %d: Failed to read audio for %s: %v", id, job, err)
			results <- Song{Title: filename, Data: nil}
			_ = cmd.Wait()
			continue
		}

		if err := cmd.Wait(); err != nil {
			log.Printf("Worker %d: yt-dlp process error for %s: %v\nStderr: %s", id, job, err, stderrBuf.String())
		}

		results <- Song{
			Title: filename,
			Data:  data,
		}
	}
}

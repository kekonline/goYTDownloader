
<img src="https://raw.githubusercontent.com/kekonline/frontend_YTDownloader/refs/heads/master/src/assets/logo.png" alt="Game Logo" width="300" />

# Go YouTube Audio Streamer

A lightweight Go web server to stream audio directly from YouTube URLs using [`yt-dlp`](https://github.com/yt-dlp/yt-dlp). It supports single audio streaming, multiple audio file downloads (zipped), and a newer streaming endpoint with timeout and size limit.

---

## Features

* Stream best audio format from a YouTube video directly.
* Support multiple URLs with concurrent downloads, packaged into a ZIP archive.
* Configurable CORS headers based on your environment.
* Dockerized for easy deployment.
* Simple and extensible codebase using Go's standard library and `yt-dlp`.
* Built-in timeout and max size limits for streams.
* Unit tests included for handler validation.

---

## Quickstart

### Prerequisites

* [Docker](https://docs.docker.com/get-docker/) installed.
* Optional: To run tests locally, install [`yt-dlp`](https://github.com/yt-dlp/yt-dlp) and have internet connectivity.
* Ensure `yt-dlp` is accessible in your PATH inside Docker container (handled by the Dockerfile).

---

### Build and Run with Docker

```bash
# Build Docker image
docker build -t my-go-audio-app .

# Run container exposing port 8080
docker run -p 8080:8080 my-go-audio-app
```

Your server will be running at `http://localhost:8080`.

---

## API Endpoints

### 1. `/api/audio-stream` (POST)

Stream the best audio of a single YouTube URL.

**Request body:**

```json
{
  "url": "https://www.youtube.com/watch?v=video_id"
}
```

**Response:**

* Streams audio in `audio/webm` format.
* Sets `Content-Disposition` with the YouTube video's filename.

---

### 2. `/api/audio-streamv2` (POST)

Improved streaming endpoint with timeout and max stream size limit (50 MB).

**Request body:**

```json
{
  "url": "https://www.youtube.com/watch?v=video_id"
}
```

---

### 3. `/api/audio-stream-multiple-files-poc` (POST)

Streams multiple URLs concurrently, then returns a ZIP archive containing all audio files.

**Request body:**

```json
{
  "url": "https://youtube.com/vid1,https://youtube.com/vid2,https://youtube.com/vid3"
}
```

* Comma-separated URLs.
* Response is a ZIP archive (`audio_files.zip`) containing each audio file.

---

## Environment Variables

The app loads `.env` for configuration. The key environment variable is:

* `HOST` — Used for CORS origin validation.

Example `.env` file:

```
HOST=http://localhost:8080
```

---

## Dockerfile Explanation

* Starts from `golang:1.24.4`.
* Installs system dependencies: `ffmpeg`, `python3-pip`, `python3-venv`.
* Creates a Python virtual environment and installs `yt-dlp`.
* Copies and downloads Go dependencies separately to leverage Docker caching.
* Copies source code, builds the Go binary.
* Exposes port 8080 and runs the built app.

---

## Running Tests

Tests focus on the HTTP handlers.

```bash
go test ./internal/handler -v
```

**Note:** Tests that stream real audio require internet access and `yt-dlp` installed on your local machine.

---

## Example Usage with `curl`

```bash
curl -X POST http://localhost:8080/api/audio-stream \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.youtube.com/watch?v=dQw4w9WgXcQ"}' --output audio.webm
```

---

## Project Structure

```
.
├── cmd/
│   └── server/          # main.go (entrypoint)
├── internal/
│   ├── handler/         # HTTP handlers for streaming
│   └── model/           # Request models (e.g., DownloadRequest)
├── Dockerfile
├── go.mod
├── go.sum
├── README.md
└── .env                 # (optional) for environment vars like HOST
```

---

## Dependencies

* [yt-dlp](https://github.com/yt-dlp/yt-dlp) (installed inside Docker and required locally for tests)
* [github.com/joho/godotenv](https://github.com/joho/godotenv) for loading `.env` files.

---

## Notes

* Audio is streamed in best available format from YouTube.
* To avoid abuse, the v2 handler limits max streamed audio size and enforces a timeout.
* The multiple files handler downloads files concurrently with worker pool.
* CORS headers allow requests only from the configured `HOST`.

FROM golang:1.24.4

RUN apt-get update && \
    apt-get install -y ffmpeg python3-pip python3-venv && \
    python3 -m venv /venv && \
    /venv/bin/pip install --upgrade pip && \
    /venv/bin/pip install yt-dlp

ENV PATH="/venv/bin:$PATH"

WORKDIR /app

# Copy go.mod and go.sum separately to leverage Docker cache for deps
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the Go app, specifying the main package directory
RUN go build -o app ./cmd/server

EXPOSE 10000

CMD ["./app"]


# go test ./internal/handler -v
# docker build -t my-go-audio-app .
# docker run -p 8080:8080 my-go-audio-app


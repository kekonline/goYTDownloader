FROM golang:1.21

# Install yt-dlp and ffmpeg
RUN apt-get update && \
    apt-get install -y ffmpeg python3-pip && \
    pip3 install yt-dlp

# Set workdir
WORKDIR /app

# Copy code
COPY . .

# Build Go app
RUN go mod tidy
RUN go build -o app

# Expose port
EXPOSE 8080

# Run binary
CMD ["./app"]

FROM golang:1.24.4

RUN apt-get update && \
    apt-get install -y ffmpeg python3-pip python3-venv && \
    python3 -m venv /venv && \
    /venv/bin/pip install --upgrade pip && \
    /venv/bin/pip install yt-dlp

ENV PATH="/venv/bin:$PATH"


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

# 1. Make changes in code

# 2. Rebuild the image
# docker build -t goytdownloader .

# 3. Stop the running container (if any)
# docker ps  # Find the container ID
# docker stop <container-id>

# 4. Run it again
# docker run -p 8080:8080 goytdownloader


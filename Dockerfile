FROM golang:1.22.2 as base

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY *.go ./

EXPOSE ${API_PORT}

# Add Air for hot reload
FROM base as development
RUN go install github.com/mitranim/gow@latest
CMD ["gow", "run", "cmd/api/main.go"]


# Command to run the executable
FROM base as production
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api/main.go
CMD ["./api"]

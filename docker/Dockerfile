# Stage 1 — build stage (kitchen)
# Full Go environment to compile the binary
FROM golang:1.26 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy dependency files first (better caching — only re-downloads if go.mod changes)
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Compile for Linux — CGO_ENABLED=0 makes a fully static binary, no external deps needed
RUN CGO_ENABLED=0 GOOS=linux go build -o raven ./cmd/broker

# Stage 2 — runtime stage (delivery box)
# Tiny Alpine image — only contains the compiled binary, nothing else
FROM alpine:latest

WORKDIR /app

# Copy just the binary from the build stage — leave the kitchen behind
COPY --from=builder /app/raven .

# Document which ports this container uses
EXPOSE 8080 9000 9001 9002

# Run the binary when the container starts
CMD ["./raven"]
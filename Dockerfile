# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go module files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o cowpoke cmd/cowpoke/main.go

# Final stage - minimal runtime image
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /build/cowpoke /cowpoke

# Run as non-root user
USER 65534:65534

ENTRYPOINT ["/cowpoke"]

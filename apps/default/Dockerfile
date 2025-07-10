FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy project files
COPY ./apps/default/* ./

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags '-extldflags "-static"' -o binary .

# Final minimal image
FROM gcr.io/distroless/static:nonroot

# Copy from builder stage
COPY --from=builder /app/binary /profile
COPY --from=builder /app/migrations /migrations

# Run as non-root user
USER 65532:65532

# Document ports
EXPOSE 80
EXPOSE 50051

# Set working directory
WORKDIR /

# Run the binary
ENTRYPOINT ["/profile"]

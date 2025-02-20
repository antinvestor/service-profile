FROM golang:1.24 AS builder

WORKDIR /

COPY go.mod .
COPY go.sum .
RUN go env -w GOFLAGS=-mod=mod && go mod download

# Copy the local package files to the container's workspace.
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o binary .

FROM gcr.io/distroless/static:nonroot

USER 65532:65532
EXPOSE 80
EXPOSE 50051

WORKDIR /


COPY --from=builder /binary /profile
COPY --from=builder /migrations /migrations

# Run the service command by default when the container starts.
ENTRYPOINT ["/profile"]

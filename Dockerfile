FROM golang:1.14 as builder

# Add Maintainer Info
LABEL maintainer="Bwire Peter <bwire517@gmail.com>"

WORKDIR /

ADD go.mod ./

RUN go mod download

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o profile_binary .

FROM scratch
COPY --from=builder /profile_binary /profile
COPY --from=builder /migrations /migrations
WORKDIR /

# Run the service command by default when the container starts.
ENTRYPOINT ["/profile"]

# Document the port that the service listens on by default.
EXPOSE 7005
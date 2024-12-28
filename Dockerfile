FROM golang:1.23 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bin/spotifyservice ./cmd/spotifyservice

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/bin/spotifyservice /app/spotifyservice

EXPOSE 8080

ENTRYPOINT ["/app/spotifyservice"]
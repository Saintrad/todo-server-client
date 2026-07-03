# ---------- Build stage ----------
FROM golang:1.25 AS builder

WORKDIR /app

ENV GOPROXY=direct
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# ---------- Run stage ----------
FROM alpine:latest

WORKDIR /app

# needed for TLS / HTTPS calls
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
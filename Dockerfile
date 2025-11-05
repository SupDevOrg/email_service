# -------- Stage 1: Build --------
FROM golang:1.23-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o main ./cmd/server


# -------- Stage 2: Distroless --------
FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=builder /app/main .
USER nonroot

EXPOSE 8080
CMD ["./main"]

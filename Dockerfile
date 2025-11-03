FROM golang:1.23-alpine

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/server

EXPOSE 8080

CMD ["./main"]
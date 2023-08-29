FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 50051

RUN go build -o identity ./cmd/identity/main.go
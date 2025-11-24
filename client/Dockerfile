FROM golang:1.25-alpine

WORKDIR /app

RUN go install github.com/githubnemo/CompileDaemon@latest

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

CMD ["CompileDaemon", "-build", "go build -o /app/main .", "-command", "/app/main"]
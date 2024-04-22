FROM golang:latest

WORKDIR /orchestrator

COPY go.mod .
COPY . .

CMD ["go", "run", "main.go"]

EXPOSE 8081
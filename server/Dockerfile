FROM golang:latest

WORKDIR /server

COPY go.mod .
COPY . .

CMD ["go", "run", "main.go"]

EXPOSE 8080
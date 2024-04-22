FROM golang:latest

WORKDIR /database

COPY go.mod .
COPY . .

CMD ["go", "run", "main.go"]

EXPOSE 5050
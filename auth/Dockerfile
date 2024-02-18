FROM golang:1.19.10-alpine3.17

WORKDIR /app

RUN ls -la

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY .. .

CMD ["go", "run", "/app/cmd/api/main.go"]

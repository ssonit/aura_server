FROM golang:1.22-alpine

RUN apk add --no-cache git curl

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
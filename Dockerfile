FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd/api/main.go

FROM scratch

# Copy the binary from the build stage
COPY --from=builder /app/main /main

# Set the command to run the binary
CMD ["/main"]
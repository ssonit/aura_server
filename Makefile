include .env

all: build

build:
	@echo "Building..."
	
	@go build -o main cmd/api/main.go

run:
	@go run cmd/api/main.go

up:
	@echo "Starting up the container..."
	docker-compose up --build -d

down:
	@echo "Stopping the container..."
	docker-compose down


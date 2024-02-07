.PHONY: run

run:
	@echo "start server..."
	go run cmd/main.go

build:
	@echo "building application..."
	go build -o bin/main ./cmd/main.go
	# go build -tags musl -o bin/main ./cmd/main.go

dev:
	@echo "starting server in development mode"
	air

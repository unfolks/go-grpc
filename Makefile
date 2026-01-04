.PHONY: run build test tidy clean

# Default target
all: build

# Run the application
run:
	go run cmd/server/main.go

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Tidy dependencies
tidy:
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs/

# Serve documentation locally (requires Python)
docs-serve:
	@echo "Serving documentation at http://localhost:8081"
	@python3 -m http.server 8081 --directory docs

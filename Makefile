# Makefile for Go tests

# Define the Go executable name
BINARY := xrf197ilz35aqbin
BIN_DIR := bin

# Default target: build the Go executable
build:
	@echo "Building the Go executable..."
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY)

# Run tests
test: build
	@echo "Running Go tests..."
	go test ./...

# Run skipped tests
test-skipped: build
	@echo "Running skipped Go tests..."
	go test -run=^TestSkipped ./...

# Clean up the built executable
clean:
	@echo "Cleaning up..."
	rm -f $(BIN_DIR)/$(BINARY)

# Help target to display available commands
help:
	@echo "Available commands:"
	@echo "  make build     - Build the Go executable"
	@echo "  make test      - Run all tests"
	@echo "  make test-skipped - Run only skipped tests"
	@echo "  make clean     - Clean up the built executable"
	@echo "  make help      - Display this help message"

TARGET ?= scanme.nmap.org
PORTS ?= 22,80,443
WORKERS ?= 5
TIMEOUT ?= 5
BINARY = portscanner

.PHONY: run build clean test fmt help

# Build the program
build:
	@echo "Building the program."
	@go build -ldflags="-s -w" -o $(BINARY) main.go

# Run the compiled program
run: build
	@echo "Running scanner with:"
	@echo "  Target: $(TARGET)"
	@echo "  Ports: $(PORTS)"
	@echo "  Workers: $(WORKERS)"
	@echo "  Timeout: $(TIMEOUT)s"
	@./$(BINARY) \
		-targets=$(TARGET) \
		-worker=$(WORKERS) \
		-ports=$(PORTS) \
		-timeout=$(TIMEOUT) \
		-booleancheck

# Format 
fmt:
	@echo "Formatting."
	@go fmt ./...

# Run tests
test:
	@echo "Running A test."
	@go test -v ./...

# Clean up
clean:
	@echo "Cleaning up."
	@rm -f $(BINARY)

# Show available commands
help:
	@echo "Makefile Commands:"
	@echo "  make build    - Build the Go program"
	@echo "  make run      - Run the compiled binary"
	@echo "  make fmt      - Format Go source code"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove compiled binaries"
	@echo "  make help     - Show available commands"

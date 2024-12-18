# Define variables
GO = go
GOFLAGS =

# Define the main target
all: build

# Build the project
build:
ifeq ($(OS),Windows_NT)
	$(GO) build $(GOFLAGS) -o life.exe
else
	$(GO) build $(GOFLAGS) -o life
endif

# Fetch dependencies
deps:
	$(GO) mod tidy

# Clean up build artifacts
clean:
ifeq ($(OS),Windows_NT)
	del /F /Q life.exe
else
	rm -f life
endif

.PHONY: all build deps clean

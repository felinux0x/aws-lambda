.PHONY: build clean deploy test

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# Binary names
BINARY_NAME=bootstrap
BINARY_DIR=bin

# Target platform for Graviton2 (ARM64)
# AWS Lambda custom runtimes on Amazon Linux 2 / 2023 expect the binary to be named "bootstrap"
GOOS=linux
GOARCH=arm64

all: test build

build:
	@echo "Building for $(GOOS)/$(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GOBUILD) -tags lambda.norpc -o $(BINARY_DIR)/$(BINARY_NAME) cmd/api/main.go

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	rm -f package.zip

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

zip: build
	@echo "Zipping binary..."
	zip -j package.zip $(BINARY_DIR)/$(BINARY_NAME)

# AWS SAM commands
sam-build:
	sam build

sam-local: sam-build
	sam local invoke ApiFunction --event events/api_event.json --env-vars env.json

deploy: sam-build
	sam deploy --guided

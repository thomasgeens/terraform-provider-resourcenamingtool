# Variables
BINARY_NAME = terraform-provider-resourcenamingtool
AUTHOR_SLUG = "thomasgeens"
REGISTRY_NAME = "registry.terraform.io"
GO_FILES = $(shell find . -type f -name '*.go')
GOBIN = "$(shell go env GOPATH)/bin"

# Default target
# Install and verify provider
verify: install
	@echo "Verifying provider installation..."; \
	cd examples/provider && terraform plan


# Install the provider for development purposes
install: build update_terraformrc
	@echo "Installing the provider..."
	@mkdir -p $(GOBIN)
	@cp $(BINARY_NAME) $(GOBIN)/$(BINARY_NAME)
	@echo "Provider installed to $(GOBIN)/$(BINARY_NAME)"

# Build the binary
build: test testacc
	@echo "Building the provider..."; \
	go build -o $(BINARY_NAME) -v ./main.go

# Update the ~/.terraformrc file using the script
update_terraformrc:
	@echo "Updating ~/.terraformrc file..."
	@./tools/update_terraformrc.sh $(REGISTRY_NAME) $(AUTHOR_SLUG) $(BINARY_NAME) $(GOBIN)

# Run unit tests
test:
	@echo "Running unit tests..."
	TF_ACC=0 go test ./internal/provider/... -v

# # Run acceptance tests (requires TF_ACC=1)
testacc:
	@echo "Running acceptance tests..."
	TF_ACC=1 go test ./internal/provider/... -v

# Run specific unit tests
test_provider:
	@echo "Running provider unit tests..."; \
	TF_ACC=0 go test ./internal/provider -run=TestProvider -v

test_datasource:
	@echo "Running datasource unit tests..."; \
	TF_ACC=0 go test ./internal/provider -run=TestAccProviderStatusDataSource -v

test_function:
	@echo "Running function unit tests..."; \
	TF_ACC=0 go test ./internal/provider -run=TestGenerateResourceNameFunction -v

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Running linters..."
	golangci-lint run

# Clean up build clutter
clean:
	@echo "Cleaning up clutter..."
	rm -f $(BINARY_NAME)
	rm -f $(GOBIN)/$(BINARY_NAME)
	rm -rf internal/provider/.terraform

# Run the provider locally
run:
	@echo "Running the provider locally..."; \
	go run main.go

# Generate documentation (if applicable)
generate: clean
	@echo "Generating documentation..."; \
    rm -rf docs/; \
	cd tools; go generate ./...

# Help
help:
	@echo "Makefile commands:"
	@echo "  verify                  - Build and verify provider installation"
	@echo "  install                 - Install the provider for development purposes"
	@echo "  build                   - Build the provider binary"
	@echo "  test                    - Run unit tests"
	@echo "  test_provider           - Run provider unit tests"
	@echo "  test_datasource         - Run datasource unit tests"
	@echo "  test_function           - Run function unit tests"
	@echo "  fmt                     - Format code"
	@echo "  lint                    - Run linters"
	@echo "  clean                   - Clean up clutter"
	@echo "  run                     - Run the provider locally"
	@echo "  generate                - Re-generate documentation"
	@echo "  update_terraformrc      - Update ~/.terraformrc file"
	@echo "  help                    - Show this help message"

.PHONY: verify install build test test_provider test_datasource test_function fmt lint clean run generate update_terraformrc help
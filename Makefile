WORKDIR         := $(shell pwd)
INTEGRATION     := memcached
BINARY_NAME      = nri-$(INTEGRATION)
GO_FILES        := ./src/
GOFLAGS          = -mod=readonly
GOLANGCI_LINT    = github.com/golangci/golangci-lint/cmd/golangci-lint
GOCOV            = github.com/axw/gocov/gocov
GOCOV_XML        = github.com/AlekSi/gocov-xml

all: build

build: clean validate test compile

clean:
	@echo "=== $(INTEGRATION) === [ clean ]: Removing binaries and coverage file..."
	@rm -rfv bin coverage.xml

validate: 
	@echo "=== $(INTEGRATION) === [ validate ]: Validating source code running gometalinter..."
	@go run $(GOFLAGS) $(GOLANGCI_LINT) run --verbose

compile: 
	@echo "=== $(INTEGRATION) === [ compile ]: Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(GO_FILES)

test: 
	@echo "=== $(INTEGRATION) === [ test ]: Running unit tests..."
	@go run $(GOFLAGS) $(GOCOV) test ./... | go run $(GOFLAGS) $(GOCOV_XML) > coverage.xml

integration-test:
	@echo "=== $(INTEGRATION) === [ test ]: running integration tests..."
	@docker-compose -f tests/docker-compose.yml pull
	@go test -v -tags=integration ./tests/. || (ret=$$?; docker-compose -f tests/docker-compose.yml down && exit $$ret)
	@docker-compose -f tests/docker-compose.yml down

# Include thematic Makefiles
include $(CURDIR)/build/ci.mk
include $(CURDIR)/build/release.mk

.PHONY: all build clean validate compile test  integration-test

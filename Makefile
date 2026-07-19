# md2docx Makefile
# Common development and build tasks.

BINARY    := md2docx
CMD_DIR   := ./cmd/md2docx
DIST_DIR  := ./dist
VERSION   ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT    ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILDDATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.buildDate=$(BUILDDATE)

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY) $(CMD_DIR)

.PHONY: build-all
build-all:
	GOOS=linux   GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-linux-amd64   $(CMD_DIR)
	GOOS=linux   GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-linux-arm64   $(CMD_DIR)
	GOOS=darwin  GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-darwin-amd64  $(CMD_DIR)
	GOOS=darwin  GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-darwin-arm64  $(CMD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-windows-amd64.exe $(CMD_DIR)

.PHONY: run
run:
	go run $(CMD_DIR)

.PHONY: run-tui
run-tui:
	go run $(CMD_DIR)

.PHONY: test
test:
	go test ./... -v -race -count=1

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: install
install:
	go install -ldflags="$(LDFLAGS)" $(CMD_DIR)

.PHONY: install-skill
install-skill:
	go run $(CMD_DIR) skill install

.PHONY: clean
clean:
	rm -rf $(DIST_DIR)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: release-dry-run
release-dry-run:
	goreleaser release --snapshot --clean --skip=publish

.DEFAULT_GOAL := build

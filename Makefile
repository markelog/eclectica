version := $(shell go run ec/main.go version)

all: install test
.PHONY: all

install:
	@echo "[+] installing dependencies"
	@go get ./...
.PHONY: install

test:
	@echo "[+] testing"
	@go test ./...
.PHONY: test

build:
	@echo "[+] building"
	@gox -osarch="darwin/amd64 linux/amd64" ./...
.PHONY: build

tag:
	@echo "[+] tagging"
	@git tag v$(version) -a -m "Release v$(version)"
.PHONY: tag

release:
	@echo "[+] releasing"
	@$(MAKE) test
	@$(MAKE) build
	@$(MAKE) tag
	@echo "[+] complete"
.PHONY: release

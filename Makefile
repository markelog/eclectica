version := $(shell go run ec/main.go version)

all: install test
.PHONY: all

install:
	@go get ./...
.PHONY: install

test:
	@go test ./...
.PHONY: test

build:
	@gox -osarch="darwin/amd64 linux/amd64" ./...
.PHONY: build

release:
	@echo "[+] releasing"
	@echo "[+] testing"
	@$(MAKE) test
	@echo "[+] building"
	@$(MAKE) build
	@echo "[+] tagging"
	@git tag v$(version) -a -m "Release v$(version)"
	@echo "[+] complete"
.PHONY: release

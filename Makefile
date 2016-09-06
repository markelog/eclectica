version := $(shell go run ec/main.go version)

install:
	@go get ./...
.PHONY: install

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

test:
	@go test ./...
.PHONY: test

build:
	@gox -osarch="darwin/amd64 linux/amd64 linux/386" ./...
.PHONY: build

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
	# We use github.com/markeog/list, which uses https://github.com/sethgrid/curse
	# which uses https://github.com/kless/term, which do not support anything
	# else then "amd64" arch – https://github.com/kless/term/issues/6
	@gox -os="darwin linux" -arch="amd64" ./...
.PHONY: build

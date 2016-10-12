
all: install test
.PHONY: all

install:
	@echo "[+] installing dependencies"
	@go get -t -v ./...
.PHONY: install

test:
	@echo "[+] testing"
	@go test -v ./...
.PHONY: test

int-test:
	$(eval tmp := $(shell pwd)"/tmp")

	@echo "[+] intergration testing"

	@rm -rf $(tmp)
	@mkdir $(tmp)

	@go build -v ./bin/ec-proxy
	@mv ec-proxy $(tmp)

	@env EC_PROXY_PLACE=$(tmp) INT=true go test -v ./bin/ec
	@rm -rf $(tmp)
.PHONY: int-test

build:
	@echo "[+] building"
	@go get github.com/mitchellh/gox
	@gox -osarch="darwin/amd64 linux/amd64" ./...
.PHONY: build

tag:
	$(eval version := $(shell go run ec/main.go version))
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

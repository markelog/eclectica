all: clean install test
.PHONY: all

clean:
	@rm -rf ec_* ec-proxy_* bin/ec/ec bin/ec/ec-proxy bin/ec-proxy/ec-proxy
.PHONY: clean

install:
	@echo "[+] installing dependencies"
	@go get -t -v ./...
.PHONY: install

test: install
	@echo "[+] testing"
	@go test -v ./...
.PHONY: test

int: install
	$(eval tmp := $(TMPDIR)"eclectica")

	@echo "[+] intergration testing"

	@rm -rf $(tmp)
	@mkdir -p $(tmp)

	@go build -v ./bin/ec-proxy
	@mv ec-proxy $(tmp)

	@env EC_PROXY_PLACE=$(tmp) EC_WITHOUT_SPINNER=true TEST_ALL=true go test -v ./bin/ec -timeout 50m

	@rm -rf $(tmp)
.PHONY: int

int-ci:
	$(eval tmp := $(TMPDIR)"eclectica")

	@echo "[+] intergration testing"

	@rm -rf $(tmp)
	@mkdir -p $(tmp)

	@go build -v ./bin/ec-proxy
	@mv ec-proxy $(tmp)

	@env EC_PROXY_PLACE=$(tmp) EC_WITHOUT_SPINNER=true go test -v ./bin/ec -timeout 50m

	@echo $(tmp)

	@rm -rf $(tmp)
	@echo $(?)
.PHONY: int-ci

build:
	@echo "[+] building"
	@go get github.com/mitchellh/gox
	@rm -rf ec_* ec-proxy_*
	@gox -osarch="darwin/amd64 linux/amd64" ./...
.PHONY: build

tag:
	$(eval version := $(shell go run bin/ec/main.go version))
	@echo "[+] tagging"
	@git tag v$(version) -a -m "Release v$(version)"
.PHONY: tag

release:
	@echo "[+] releasing"
	@$(MAKE) clean
	@$(MAKE) test
	@$(MAKE) build
	@$(MAKE) tag
	@echo "[+] complete"
.PHONY: release

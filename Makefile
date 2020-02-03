GO_SOURCES = $(shell find . -type f -name '*.go' ! -path '**/mocks/*')

ifeq ($(OS),Windows_NT)
	OUTPUT=streaming-http-adapter.exe
else
	OUTPUT=streaming-http-adapter
endif

ifeq (,$(shell go env GOBIN))
	GOBIN=$(shell go env GOPATH)/bin
else
	GOBIN=$(shell go env GOBIN)
endif

.PHONY: all
all: build test

.PHONY: build
build: $(OUTPUT) ## Build the executable for current architecture (local dev)

$(OUTPUT): $(GO_SOURCES)
	go build -o $(OUTPUT) -gcflags="all=-N -l" main.go

pkg/rpc/riff-rpc.pb.go: riff-rpc.proto
	protoc -I . riff-rpc.proto --go_out=plugins=grpc:pkg/rpc

.PHONY: release
release: verify-mocks test streaming-http-adapter-linux-amd64.tgz ## Build the executable as a static linux executable

streaming-http-adapter-linux-amd64.tgz: $(GO_SOURCES)
	mkdir temp \
	&& CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o temp/streaming-http-adapter main.go \
	&& tar -czf streaming-http-adapter-linux-amd64.tgz -C temp/ streaming-http-adapter \
	&& rm -fR temp

.PHONY: clean
clean: ## Clean generated files
	rm -f $(OUTPUT)
	rm -f streaming-http-adapter-linux-amd64.tgz

.PHONY: test
test: ## Run the tests
	go test ./...

pkg/proxy/mocks/RiffClient.go: pkg/rpc/riff-rpc.pb.go
	$(MOCKERY) -output ./pkg/proxy/mocks -dir ./pkg/rpc -name RiffClient

pkg/proxy/mocks/Riff_InvokeClient.go: pkg/rpc/riff-rpc.pb.go
	$(MOCKERY) -output ./pkg/proxy/mocks -dir ./pkg/rpc -name Riff_InvokeClient

.PHONY: gen-mocks
gen-mocks: pkg/proxy/mocks/RiffClient.go pkg/proxy/mocks/Riff_InvokeClient.go

.PHONY: clean-mocks
clean-mocks: ## Delete mocks
	rm -fR pkg/proxy/mocks

.PHONY: verify-mocks
verify-mocks: ## Verify that mocks are up to date
	$(MOCKERY) -print -dir ./pkg/rpc -name RiffClient | diff ./pkg/proxy/mocks/RiffClient.go  -
	$(MOCKERY) -print -dir ./pkg/rpc -name Riff_InvokeClient | diff ./pkg/proxy/mocks/Riff_InvokeClient.go  -
	
.PHONY: install-mockery
install-mockery: ## Install mockery if missing
ifeq (, $(shell which mockery))
	@# Use go get in GOPATH mode to install/update mockery. This avoids polluting go.mod/go.sum.
	GO111MODULE=off go get github.com/vektra/mockery/.../
	MOCKERY=$(GOBIN)/mockery
else
	MOCKERY=$(shell which mockery)
endif


# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

GO_SOURCES = $(shell find . -type f -name '*.go')
ifeq ($(OS),Windows_NT)
	OUTPUT=streaming-http-adapter.exe
else
	OUTPUT=streaming-http-adapter
endif

.PHONY: build
build: $(OUTPUT) ## Build the executable for current architecture (local dev)

$(OUTPUT): $(GO_SOURCES)
	go build -o $(OUTPUT) main.go

pkg/rpc/riff-rpc.pb.go: riff-rpc.proto
	protoc -I . riff-rpc.proto --go_out=plugins=grpc:pkg/rpc

.PHONY: release
release: $(GO_SOURCES) ## Build the executable as a static linux executable
	CGO_ENABLED=0 GOOS=linux go build \
		-o $(OUTPUT)-linux main.go

.PHONY: clean
clean: ## Clean generated files
	rm -f $(OUTPUT)

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

LDFLAGS = -ldflags="-X 'main.mode=release' -X 'github.com/calico32/goose/lib/std/platform.BuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')'"

.PHONY: builddir goose goose-docs goose-lsp goose-vsix goose-darwin-arm64 goose.wasm install all platforms goose-windows-amd64 goose-linux-amd64 clean test
goose: builddir
	go build $(LDFLAGS) -o build/goose ./bin/cmd/goose.go

goose-docs: builddir
	go generate ./...
	go build -tags release $(LDFLAGS) -o build/goose-docs ./bin/docs/

goose-lsp: builddir
	go build $(LDFLAGS) -o build/goose-lsp ./bin/lsp/main.go

goose-vsix: builddir
	cd ./extension; bun run build
	cp ./extension/goose-* ./build

goose.wasm: builddir
	tinygo build -o build/goose.wasm wasm/main.go
# GOOS=js GOARCH=wasm go build -o goose.wasm wasm/main.go

all: goose goose-docs goose-lsp goose-vsix goose.wasm

install: goose goose-lsp
	cp build/goose $(HOME)/.local/bin
	cp build/goose-lsp $(HOME)/.local/bin

builddir:
	mkdir -p build

# Cross compilation
platforms: goose-windows-amd64 goose-linux-amd64 goose-darwin-arm64

goose-windows-amd64: builddir
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o build/goose-windows-amd64.exe ./bin/cmd/goose.go

goose-linux-amd64: builddir
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/goose-linux-amd64 ./bin/cmd/goose.go

goose-darwin-arm64: builddir
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o build/goose-darwin-arm64 ./bin/cmd/goose.go

# Clean
clean:
	rm -rf build

# Run tests
test:
	go test -v ./...

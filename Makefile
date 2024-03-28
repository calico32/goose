LDFLAGS = -X 'main.mode=release'

.PHONY: goose goose-darwin-arm64 install all goose.wasm
goose:
	go build -ldflags="$(LDFLAGS) -X 'github.com/calico32/goose/lib/std/platform.BuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')'" -o goose ./bin/cmd/goose.go

goose-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o goose-darwin-arm64 ./bin/cmd/goose.go

install: goose
	cp goose $(HOME)/.local/bin

goose.wasm: wasm.target.json
	tinygo build -o goose.wasm -target ./wasm.target.json -wasm-abi js wasm/main.go
# GOOS=js GOARCH=wasm go build -o goose.wasm wasm/main.go

all: goose goose.wasm

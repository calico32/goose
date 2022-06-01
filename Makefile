CFLAGS = '-X main.mode=release'

goose: **/*.go
	go build -o goose cli/main.go

goose.wasm: **/*.go
	GOOS=js GOARCH=wasm go build -o goose.wasm wasm/main.go

all: goose goose.wasm

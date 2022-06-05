CFLAGS = '-X main.mode=release'

goose: **/*.go
	go build -o goose cli/goose.go

goose.wasm: **/*.go
	tinygo build -o goose.wasm -target wasm -wasm-abi js wasm/main.go

all: goose goose.wasm

CFLAGS = '-X main.mode=release'

goose: **/*.go
	go build -o goose cli/goose.go

goose.wasm: **/*.go wasm.target.json
	tinygo build -o goose.wasm -target ./wasm.target.json -wasm-abi js wasm/main.go

all: goose goose.wasm

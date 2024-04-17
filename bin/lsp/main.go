package main

import (
	"context"

	"github.com/calico32/goose/lsp"
)

func main() {
	ctx := context.Background()
	err := lsp.RunServerOnAddress(ctx, "localhost:4389")
	if err != nil {
		panic(err)
	}
}

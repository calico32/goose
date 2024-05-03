package main

import (
	"context"

	"os"

	"github.com/calico32/goose/lsp"
	"go.lsp.dev/jsonrpc2"
)

func main() {
	ctx := context.Background()
	transport := &stdioTransport{}
	stream := jsonrpc2.NewStream(transport)
	_, err := lsp.Start(ctx, stream)
	if err != nil {
		panic(err)
	}
}

type stdioTransport struct{}

func (t *stdioTransport) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (t *stdioTransport) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (t *stdioTransport) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}

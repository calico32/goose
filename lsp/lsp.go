package lsp

import (
	"context"
	"net"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func RunServerOnAddress(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		_, err = New(ctx, jsonrpc2.NewStream(conn))
		if err != nil {
			return err
		}
	}
}

func New(ctx context.Context, stream jsonrpc2.Stream) (*LanguageServer, error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	s := &LanguageServer{
		logger: log,
	}

	// jsonrpcOpts := []jsonrpc2.Options{
	// 	jsonrpc2.WithCanceler(protocol.Canceller),
	// 	jsonrpc2.WithCapacity(protocol.DefaultBufferSize),
	// 	jsonrpc2.WithLogger(s.logger.Named("jsonrpc2")),
	// }
	_, s.conn, s.client = protocol.NewServer(ctx, s, stream, zap.NewNop())

	logger := s.logger.Named("server")
	WithContext(ctx, logger)

	return s, nil
}

type LanguageServer struct {
	conn   jsonrpc2.Conn
	client protocol.Client
	logger *zap.Logger
}

type ctxKey struct{}

// WithContext returns a context.Context with *zap.Logger as a context value.
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// FromContext extracts *zap.Logger from a given context.Context and returns it.
func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ctxKey{}).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}
	return logger
}

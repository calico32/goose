package lsp

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func nopTimeEncoder(_ time.Time, _ zapcore.PrimitiveArrayEncoder) {}

func Start(ctx context.Context, stream jsonrpc2.Stream) (*LanguageServer, error) {
	id := "lsp+" + fmt.Sprint(rand.Int63())
	s := &LanguageServer{
		id:          id,
		state:       ServerStateIdle,
		sourceFiles: make(map[protocol.DocumentURI]*Mutexed[[]byte]),
		modules:     make(map[protocol.DocumentURI]*Mutexed[*ast.Module]),
		uris:        make(map[string]protocol.DocumentURI),
		fsets:       make(map[protocol.DocumentURI]*token.FileSet),
		parseErrors: make(map[protocol.DocumentURI][]protocol.Diagnostic),
	}

	zap.RegisterSink(id, func(url *url.URL) (zap.Sink, error) {
		s := &lspSink{server: s, ctx: ctx}
		return s, nil
	})

	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zap.DebugLevel)
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeTime = nopTimeEncoder
	// config.OutputPaths = append(config.OutputPaths, "file:/tmp/goose-lsp.log")
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	s.logger = logger
	s.traceWriter = newTraceWriter(s.logger)
	s.timer = newTimer(s.logger)

	// jsonrpcOpts := []jsonrpc2.Options{
	// 	jsonrpc2.WithCanceler(protocol.Canceller),
	// 	jsonrpc2.WithCapacity(protocol.DefaultBufferSize),
	// 	jsonrpc2.WithLogger(s.logger.Named("jsonrpc2")),
	// }
	_, s.conn, s.client = protocol.NewServer(ctx, s, stream, zap.NewNop())

	logger.Sugar().Infof("goose language server started (pid: %d)", os.Getpid())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		logger.Sugar().Info("shutting down")
		s.conn.Close()
		os.Exit(0)
	}()

	<-ctx.Done()

	return s, nil
}

type Mutexed[T any] struct {
	mu sync.Mutex
	v  T
}

func (m *Mutexed[T]) Lock() T {
	// m.mu.Lock()
	return m.v
}

func (m *Mutexed[T]) Update(v T) {
	m.v = v
	// m.mu.Unlock()
}

func (m *Mutexed[T]) Unlock() {
	// m.mu.Unlock()
}

type LanguageServer struct {
	id          string
	conn        jsonrpc2.Conn
	client      protocol.Client
	logger      *zap.Logger
	state       ServerState
	traceWriter io.Writer
	timer       *timer

	fsets       map[protocol.DocumentURI]*token.FileSet
	sourceFiles map[protocol.DocumentURI]*Mutexed[[]byte]
	modules     map[protocol.DocumentURI]*Mutexed[*ast.Module]
	uris        map[string]protocol.DocumentURI
	parseErrors map[protocol.DocumentURI][]protocol.Diagnostic
}

type ServerState int

const (
	ServerStateIdle ServerState = iota
	ServerStateRunning
	ServerStateShuttingDown
	ServerStateStopped
)

type lspSink struct {
	server *LanguageServer
	ctx    context.Context
}

func (s *lspSink) Write(p []byte) (n int, err error) {
	err = s.server.client.LogMessage(s.ctx, &protocol.LogMessageParams{
		Type:    protocol.MessageTypeLog,
		Message: string(p),
	})

	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (s *lspSink) Sync() error {
	return nil
}

func (s *lspSink) Close() error {
	return nil
}

type traceWriter struct {
	logger *zap.Logger
	buffer strings.Builder
}

func newTraceWriter(logger *zap.Logger) *traceWriter {
	return &traceWriter{logger: logger.WithOptions(zap.AddCallerSkip(4))}
}

func (w *traceWriter) Write(p []byte) (n int, err error) {
	w.buffer.Write(p)
	w.flush()
	return len(p), nil
}

func (w *traceWriter) flush() {
	s := w.buffer.String()
	for {
		firstNewline := strings.IndexByte(s, '\n')
		if firstNewline == -1 {
			break
		}

		line := s[:firstNewline]
		w.logger.Debug(line)
		s = s[firstNewline+1:]
	}
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

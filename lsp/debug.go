package lsp

import (
	"time"

	"go.uber.org/zap"
)

type timer struct {
	logger *zap.SugaredLogger
}

type timerMark struct {
	*timer
	name  string
	start time.Time
	done  bool
}

func newTimer(logger *zap.Logger) *timer {
	return &timer{logger: logger.WithOptions(zap.AddCallerSkip(1)).Sugar()}
}

func (t *timer) Mark(name string) *timerMark {
	t.logger.Debugf("[timer] %s started", name)
	return &timerMark{
		timer: t,
		name:  name,
		start: time.Now(),
	}
}

func (m *timerMark) Done() {
	if m.done {
		return
	}
	m.done = true
	elapsed := time.Since(m.start)
	m.logger.Debugf("[timer] %s: %s", m.name, elapsed)
}

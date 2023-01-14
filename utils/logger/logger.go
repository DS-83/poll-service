package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

const correlationID string = "Correlation-ID"

type Logger struct {
	logger *logrus.Logger
}

func NewLogger(l *logrus.Logger) *Logger {
	return &Logger{
		logger: l,
	}
}

func (l *Logger) Debug(c context.Context, args ...interface{}) {
	id := c.Value(correlationID)
	entry := l.logger.WithField(correlationID, id)
	entry.Debug(args...)
}

func (l *Logger) Error(c context.Context, args ...interface{}) {
	id := c.Value(correlationID)
	entry := l.logger.WithField(correlationID, id)
	entry.Error(args...)
}

func (l *Logger) Info(c context.Context, args ...interface{}) {
	id := c.Value(correlationID)
	entry := l.logger.WithField(correlationID, id)
	entry.Info(args...)
}

// Package logger defined logger that will be used by tasks
package logger

import (
	"log/slog"
	"os"

	"github.com/google/uuid"
)

type logHandler struct {
	logger     *slog.Logger
	attributes []slog.Attr
}

func (lh *logHandler) AddAttribute(attrs []slog.Attr) {
	lh.attributes = append(lh.attributes, attrs...)
}

func (lh *logHandler) Log(message string) {
	anyAttributes := make([]any, 0, len(lh.attributes))
	for _, v := range lh.attributes {
		anyAttributes = append(anyAttributes, v)
	}
	lh.logger.Info(message, anyAttributes...)
}

func NewLogHandler(logger *slog.Logger) logHandler {
	attributes := []slog.Attr{
		slog.String("log_id", uuid.NewString()),
	}
	logHandler := logHandler{logger: logger, attributes: attributes}
	return logHandler
}

func NewJSONLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return slog.New(handler)
}

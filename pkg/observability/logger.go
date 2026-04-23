package observability

import (
	"context"
	"log/slog"

	"github.com/aws/aws-xray-sdk-go/xray"
)

// LoggerWithTrace returns a logger decorated with the AWS X-Ray Trace ID if it exists in the context.
// This allows for correlation between traces and logs.
func LoggerWithTrace(ctx context.Context, logger *slog.Logger) *slog.Logger {
	traceID := xray.TraceID(ctx)
	if traceID != "" {
		return logger.With("trace_id", traceID)
	}
	return logger
}

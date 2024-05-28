package middleware

import (
	"encoding/json"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("sentry")
)

type LogData struct {
	TraceID    string                 `json:"trace_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	Attributes map[string]interface{} `json:"attributes"`
}

func (ld *LogData) ToJSON() (string, error) {
	bytes, err := json.Marshal(ld)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func LoggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Start tracing
	ctx, span := tracer.Start(c.Context(), "request")
	defer span.End()

	// Pass the context with the span
	c.SetUserContext(ctx)

	// Call next handler
	err := c.Next()

	// Capture logs and errors
	elapsed := time.Since(start)
	logData := LogData{
		TraceID:   span.SpanContext().TraceID().String(),
		Timestamp: start,
		Level:     "info",
		Message:   "Request completed",
		Attributes: map[string]interface{}{
			"method":   c.Method(),
			"endpoint": c.Path(),
			"status":   c.Response().StatusCode(),
			"elapsed":  elapsed.String(),
		},
	}
	if err != nil {
		logData.Level = "error"
		logData.Message = err.Error()

		// Send error to Sentry
		sentry.CaptureException(err)
	}

	// Convert log to JSON
	logJSON, jsonErr := logData.ToJSON()
	if jsonErr != nil {
		log.Printf("Failed to marshal log data: %v", jsonErr)
		return err
	}

	// Output the log
	log.Println(logJSON)

	return err
}

package logger

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *LoggerConfig

type LoggerConfig struct {
	ServiceName string
	Logger      *logrus.Logger
}

type LogMessage struct {
	File    string
	Level   string
	Message string
	TraceId string
}

const _LogMessageKey = "log_message_key"

func NewLoggerConfig(serviceName string) *LoggerConfig {
	logger := initializeLogger()
	Log = &LoggerConfig{
		ServiceName: serviceName,
		Logger:      logger,
	}
	return Log
}

// initializeLogger initializes a new logger instance.
func initializeLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	return log
}

// Middleware logs the details of each request.
func (lc *LoggerConfig) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		// Create a new context with a sync.Map for storing log messages
		ctx := context.WithValue(c.UserContext(), _LogMessageKey, &sync.Map{})
		c.SetUserContext(ctx)

		err := c.Next()

		// Get trace ID from context
		traceID := c.Locals("trace_id")
		statusCode := c.Response().StatusCode()

		logEntry := lc.Logger.WithFields(logrus.Fields{
			"method":   c.Method(),
			"url":      c.OriginalURL(),
			"duration": time.Since(start),
			"status":   statusCode,
			"trace_id": traceID,
		})

		if err != nil {
			logEntry = logEntry.WithField("error", err.Error())
			logEntry.Error("Request failed")
		} else {
			logEntry.Info("Request completed")
		}

		// Retrieve and log accumulated log messages
		value, ok := extract(ctx)
		if ok {
			if tmp, ok := value.Load(_LogMessageKey); ok {
				if lms, ok := tmp.([]LogMessage); ok {
					for _, lm := range lms {
						lc.Logger.WithFields(logrus.Fields{
							"file":     lm.File,
							"level":    lm.Level,
							"message":  lm.Message,
							"trace_id": lm.TraceId,
						}).Error("Logged message")
					}
				}
			}
		}

		return err
	}
}

func (lc *LoggerConfig) Errorf(ctx context.Context, format string, args ...interface{}) {
	var (
		lms  []LogMessage
		file string
	)

	value, ok := extract(ctx)
	if !ok {
		value = &sync.Map{}
		ctx = context.WithValue(ctx, _LogMessageKey, value)
	}

	_, fileName, line, _ := runtime.Caller(1)
	file = fmt.Sprintf("%s:%d", fileName, line)

	files := strings.Split(fileName, "/")
	if len(files) > 4 {
		file = fmt.Sprintf("%s:%d", strings.Join(files[len(files)-3:], "/"), line)
	}

	tmp, ok := value.LoadAndDelete(_LogMessageKey)
	if ok {
		lms = tmp.([]LogMessage)
	}

	// Extract trace ID from context using OpenTelemetry's built-in function
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()

	lms = append(lms, LogMessage{
		File:    file,
		Level:   logrus.ErrorLevel.String(),
		Message: fmt.Sprintf(format, args...),
		TraceId: traceID, // Assuming LogMessage has a TraceID field
	})

	value.Store(_LogMessageKey, lms)

}

func extract(ctx context.Context) (*sync.Map, bool) {
	value, ok := ctx.Value(_LogMessageKey).(*sync.Map)
	return value, ok
}

func (lc *LoggerConfig) Infof(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	lc.Logger.WithFields(logrus.Fields{
		"trace_id": ctx.Value("trace_id"), // Assuming trace_id is stored in the context
	}).Info(message)
}

func (l *LoggerConfig) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *LoggerConfig) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *LoggerConfig) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

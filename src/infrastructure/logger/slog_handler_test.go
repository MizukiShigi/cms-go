package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"testing"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
)

func TestNewHandler(t *testing.T) {
	var buf bytes.Buffer
	jsonHandler := slog.NewJSONHandler(&buf, nil)
	
	handler := NewHandler(jsonHandler)
	
	if handler == nil {
		t.Error("NewHandler() should return a handler")
	}
	
	// Check that it's our custom Handler type
	customHandler, ok := handler.(Handler)
	if !ok {
		t.Error("NewHandler() should return Handler type")
	}
	
	if customHandler.handler != jsonHandler {
		t.Error("NewHandler() should wrap the provided handler")
	}
}

func TestHandler_Handle_WithLoggingContext(t *testing.T) {
	var buf bytes.Buffer
	jsonHandler := slog.NewJSONHandler(&buf, nil)
	handler := NewHandler(jsonHandler)
	
	// Create context with logging values
	loggingMap := &sync.Map{}
	loggingMap.Store("user_id", "test-user-123")
	loggingMap.Store("request_id", "req-456")
	loggingMap.Store("operation", "test_operation")
	
	ctx := context.WithValue(context.Background(), domaincontext.Logging, loggingMap)
	
	// Create a log record
	record := slog.NewRecord(slog.TimeNow(), slog.LevelInfo, "Test message", 0)
	
	err := handler.Handle(ctx, record)
	if err != nil {
		t.Errorf("Handle() error = %v", err)
	}
	
	// Parse the JSON output
	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Errorf("Failed to parse log output: %v", err)
	}
	
	// Check that logging context values were added
	if logEntry["user_id"] != "test-user-123" {
		t.Errorf("Expected user_id in log output, got %v", logEntry["user_id"])
	}
	if logEntry["request_id"] != "req-456" {
		t.Errorf("Expected request_id in log output, got %v", logEntry["request_id"])
	}
	if logEntry["operation"] != "test_operation" {
		t.Errorf("Expected operation in log output, got %v", logEntry["operation"])
	}
	
	// Check that the original message is preserved
	if logEntry["msg"] != "Test message" {
		t.Errorf("Expected original message in log output, got %v", logEntry["msg"])
	}
}

func TestHandler_Handle_WithoutLoggingContext(t *testing.T) {
	var buf bytes.Buffer
	jsonHandler := slog.NewJSONHandler(&buf, nil)
	handler := NewHandler(jsonHandler)
	
	// Create context without logging values
	ctx := context.Background()
	
	// Create a log record
	record := slog.NewRecord(slog.TimeNow(), slog.LevelInfo, "Test message", 0)
	
	err := handler.Handle(ctx, record)
	if err != nil {
		t.Errorf("Handle() error = %v", err)
	}
	
	// Parse the JSON output
	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Errorf("Failed to parse log output: %v", err)
	}
	
	// Check that the original message is preserved
	if logEntry["msg"] != "Test message" {
		t.Errorf("Expected original message in log output, got %v", logEntry["msg"])
	}
	
	// Should not have additional logging context fields
	if _, exists := logEntry["user_id"]; exists {
		t.Error("Should not have user_id without logging context")
	}
}

func TestHandler_Handle_WithNonStringKeys(t *testing.T) {
	var buf bytes.Buffer
	jsonHandler := slog.NewJSONHandler(&buf, nil)
	handler := NewHandler(jsonHandler)
	
	// Create context with non-string keys (should be ignored)
	loggingMap := &sync.Map{}
	loggingMap.Store("valid_key", "valid_value")
	loggingMap.Store(123, "invalid_key_value") // Non-string key
	
	ctx := context.WithValue(context.Background(), domaincontext.Logging, loggingMap)
	
	// Create a log record
	record := slog.NewRecord(slog.TimeNow(), slog.LevelInfo, "Test message", 0)
	
	err := handler.Handle(ctx, record)
	if err != nil {
		t.Errorf("Handle() error = %v", err)
	}
	
	// Parse the JSON output
	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Errorf("Failed to parse log output: %v", err)
	}
	
	// Check that valid key is included
	if logEntry["valid_key"] != "valid_value" {
		t.Errorf("Expected valid_key in log output, got %v", logEntry["valid_key"])
	}
	
	// Check that invalid key is not included
	if _, exists := logEntry["123"]; exists {
		t.Error("Non-string keys should be ignored")
	}
}

func TestHandler_Enabled(t *testing.T) {
	var buf bytes.Buffer
	jsonHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})
	handler := NewHandler(jsonHandler)
	
	ctx := context.Background()
	
	tests := []struct {
		level   slog.Level
		enabled bool
	}{
		{slog.LevelDebug, false},
		{slog.LevelInfo, false},
		{slog.LevelWarn, true},
		{slog.LevelError, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			enabled := handler.Enabled(ctx, tt.level)
			if enabled != tt.enabled {
				t.Errorf("Enabled(%v) = %v, want %v", tt.level, enabled, tt.enabled)
			}
		})
	}
}

func TestHandler_WithAttrs(t *testing.T) {
	var buf bytes.Buffer
	jsonHandler := slog.NewJSONHandler(&buf, nil)
	handler := NewHandler(jsonHandler)
	
	attrs := []slog.Attr{
		slog.String("component", "test"),
		slog.Int("version", 1),
	}
	
	newHandler := handler.WithAttrs(attrs)
	
	// Check that it returns our Handler type
	customHandler, ok := newHandler.(Handler)
	if !ok {
		t.Error("WithAttrs() should return Handler type")
	}
	
	// The wrapped handler should have the attributes
	if customHandler.handler == handler.(Handler).handler {
		t.Error("WithAttrs() should create a new handler with attributes")
	}
}

func TestHandler_WithGroup(t *testing.T) {
	var buf bytes.Buffer
	jsonHandler := slog.NewJSONHandler(&buf, nil)
	handler := NewHandler(jsonHandler)
	
	groupName := "test_group"
	newHandler := handler.WithGroup(groupName)
	
	// Check that it returns our Handler type
	customHandler, ok := newHandler.(Handler)
	if !ok {
		t.Error("WithGroup() should return Handler type")
	}
	
	// The wrapped handler should have the group
	if customHandler.handler == handler.(Handler).handler {
		t.Error("WithGroup() should create a new handler with group")
	}
}
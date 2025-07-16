package logger

import (
	"testing"
)

func TestNewLogger_NotNil(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Error("Expected logger to be non-nil")
	}
}

func TestNewLogger_Singleton(t *testing.T) {
	logger1 := NewLogger()
	logger2 := NewLogger()
	if logger1 != logger2 {
		t.Error("Expected NewLogger to return the same instance (singleton)")
	}
}

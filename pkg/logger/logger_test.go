package logger

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// captureOutput captures the output of the logger
func captureOutput(f func()) string {
	// Save the original output
	originalOutput := os.Stdout

	// Create a pipe to capture the output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Reset the logger to use the pipe
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log = zerolog.New(w).With().Timestamp().Logger()

	// Call the function that produces output
	f()

	// Close the writer and restore the original output
	w.Close()
	os.Stdout = originalOutput

	// Read the output from the pipe
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func TestInit(t *testing.T) {
	// Test with valid log level
	Init("debug", false)
	if zerolog.GlobalLevel() != zerolog.DebugLevel {
		t.Errorf("Expected log level to be debug, got %v", zerolog.GlobalLevel())
	}

	// Test with invalid log level (should default to info)
	Init("invalid", false)
	if zerolog.GlobalLevel() != zerolog.InfoLevel {
		t.Errorf("Expected log level to be info, got %v", zerolog.GlobalLevel())
	}

	// Test with pretty output
	Init("info", true)
	// Can't easily test the pretty output, but at least ensure it doesn't crash
}

func TestDebug(t *testing.T) {
	output := captureOutput(func() {
		Debug("test debug message")
	})

	if !strings.Contains(output, "test debug message") {
		t.Errorf("Expected output to contain 'test debug message', got %s", output)
	}

	if !strings.Contains(output, "debug") {
		t.Errorf("Expected output to contain 'debug', got %s", output)
	}
}

func TestInfo(t *testing.T) {
	output := captureOutput(func() {
		Info("test info message")
	})

	if !strings.Contains(output, "test info message") {
		t.Errorf("Expected output to contain 'test info message', got %s", output)
	}

	if !strings.Contains(output, "info") {
		t.Errorf("Expected output to contain 'info', got %s", output)
	}
}

func TestWarn(t *testing.T) {
	output := captureOutput(func() {
		Warn("test warn message")
	})

	if !strings.Contains(output, "test warn message") {
		t.Errorf("Expected output to contain 'test warn message', got %s", output)
	}

	if !strings.Contains(output, "warn") {
		t.Errorf("Expected output to contain 'warn', got %s", output)
	}
}

func TestError(t *testing.T) {
	err := errors.New("test error")
	output := captureOutput(func() {
		Error(err, "test error message")
	})

	if !strings.Contains(output, "test error message") {
		t.Errorf("Expected output to contain 'test error message', got %s", output)
	}

	if !strings.Contains(output, "error") {
		t.Errorf("Expected output to contain 'error', got %s", output)
	}

	if !strings.Contains(output, "test error") {
		t.Errorf("Expected output to contain 'test error', got %s", output)
	}
}

func TestInfof(t *testing.T) {
	output := captureOutput(func() {
		Infof("test %s message", "formatted")
	})

	if !strings.Contains(output, "test formatted message") {
		t.Errorf("Expected output to contain 'test formatted message', got %s", output)
	}

	if !strings.Contains(output, "info") {
		t.Errorf("Expected output to contain 'info', got %s", output)
	}
}

func TestErrorf(t *testing.T) {
	err := errors.New("test error")
	output := captureOutput(func() {
		Errorf(err, "test %s message", "formatted")
	})

	if !strings.Contains(output, "test formatted message") {
		t.Errorf("Expected output to contain 'test formatted message', got %s", output)
	}

	if !strings.Contains(output, "error") {
		t.Errorf("Expected output to contain 'error', got %s", output)
	}

	if !strings.Contains(output, "test error") {
		t.Errorf("Expected output to contain 'test error', got %s", output)
	}
}

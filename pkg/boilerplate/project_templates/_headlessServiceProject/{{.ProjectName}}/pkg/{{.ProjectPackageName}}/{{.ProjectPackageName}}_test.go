package {{.ProjectPackageName}}

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestNewService(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Test with default interval (0)
	service := NewService(logger, 0)
	require.NotNil(t, service)
	assert.Equal(t, logger, service.logger)
	assert.Equal(t, 5*time.Second, service.interval) // Should default to 5 seconds

	// Test with custom interval
	customInterval := 2 * time.Second
	service2 := NewService(logger, customInterval)
	require.NotNil(t, service2)
	assert.Equal(t, customInterval, service2.interval)
}

func TestServiceStartStop(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Use very short interval for testing
	service := NewService(logger, 100*time.Millisecond)

	// Start service in background
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- service.Start(ctx)
	}()

	// Give service a moment to start
	time.Sleep(50 * time.Millisecond)

	// Service should be running and logging
	// We can't easily test the HTTP server without complex setup,
	// but we can test that Start() returns when context is cancelled

	// Wait for service to stop when context times out
	select {
	case serviceErr := <-done:
		require.NoError(t, serviceErr)
	case <-time.After(1 * time.Second):
		t.Fatal("Service did not stop within timeout")
	}
}

func TestServiceStartCancellation(t *testing.T) {
	logger := zaptest.NewLogger(t)
	service := NewService(logger, 10*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- service.Start(ctx)
	}()

	// Give service a moment to start
	time.Sleep(50 * time.Millisecond)

	// Cancel context to trigger shutdown
	cancel()

	// Wait for service to stop
	select {
	case serviceErr := <-done:
		require.NoError(t, serviceErr)
	case <-time.After(2 * time.Second):
		t.Fatal("Service did not stop within timeout")
	}
}
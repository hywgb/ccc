package esl

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

// These tests require a running FreeSWITCH instance with ESL on port 8021.
// Set ESL_INTEGRATION_TEST=1 to run.
func testClient(t *testing.T) *Client {
	if os.Getenv("ESL_INTEGRATION_TEST") != "1" {
		t.Skip("Set ESL_INTEGRATION_TEST=1 to run ESL integration tests")
	}

	host := os.Getenv("ESL_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	pass := os.Getenv("ESL_PASSWORD")
	if pass == "" {
		pass = "ClueCon"
	}

	logger := zerolog.New(zerolog.NewTestWriter(t))
	c := NewClient(Config{
		Host:     host,
		Port:     8021,
		Password: pass,
		PoolSize: 2,
		Logger:   logger,
	})
	return c
}

func TestIntegration_Connect(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test basic API command
	result, err := c.SendCommand(ctx, "status")
	if err != nil {
		t.Fatalf("SendCommand(status) failed: %v", err)
	}
	if !strings.Contains(result, "FreeSWITCH") {
		t.Errorf("Expected FreeSWITCH in status output, got: %s", result)
	}
	t.Logf("FreeSWITCH status: %s", result)
}

func TestIntegration_ConnectionPool(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Acquire multiple connections to test pool
	cn1, err := c.Acquire(ctx)
	if err != nil {
		t.Fatalf("Acquire #1 failed: %v", err)
	}
	cn2, err := c.Acquire(ctx)
	if err != nil {
		t.Fatalf("Acquire #2 failed: %v", err)
	}

	if cn1.id == cn2.id {
		t.Error("Pool returned same connection twice")
	}

	c.Release(cn1)
	c.Release(cn2)

	// Re-acquire should reuse existing connections
	cn3, err := c.Acquire(ctx)
	if err != nil {
		t.Fatalf("Acquire #3 failed: %v", err)
	}
	if !cn3.connected {
		t.Error("Re-acquired connection should be connected")
	}
	c.Release(cn3)
}

func TestIntegration_Originate_InvalidDest(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Originate to invalid destination should return error
	_, err := c.Originate(ctx, "user/nonexistent", "1000", "&park()")
	if err == nil {
		t.Error("Expected error for invalid originate destination")
	}
	t.Logf("Originate error (expected): %v", err)
}

func TestIntegration_HangupCall_InvalidUUID(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Hangup with invalid UUID should return error
	err := c.HangupCall(ctx, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("Expected error for invalid UUID hangup")
	}
	t.Logf("Hangup error (expected): %v", err)
}

func TestIntegration_HoldRetrieve_InvalidUUID(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// NOTE: FreeSWITCH uuid_hold does NOT return -ERR for non-existent UUIDs.
	// This is expected FS behavior — the command silently succeeds.
	err := c.HoldCall(ctx, "00000000-0000-0000-0000-000000000000")
	t.Logf("Hold result (FS returns +OK even for invalid UUID): err=%v", err)

	err = c.RetrieveCall(ctx, "00000000-0000-0000-0000-000000000000")
	t.Logf("Retrieve result: err=%v", err)
}

func TestIntegration_SendDTMF_InvalidUUID(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.SendDTMF(ctx, "00000000-0000-0000-0000-000000000000", "123")
	if err == nil {
		t.Error("Expected error for invalid UUID DTMF")
	}
}

func TestIntegration_Transfer_InvalidUUID(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.TransferCall(ctx, "00000000-0000-0000-0000-000000000000", "9999")
	if err == nil {
		t.Error("Expected error for invalid UUID transfer")
	}
}

func TestIntegration_Bridge_InvalidUUIDs(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.Bridge(ctx, "00000000-0000-0000-0000-000000000000", "11111111-1111-1111-1111-111111111111")
	if err == nil {
		t.Error("Expected error for invalid UUID bridge")
	}
}

func TestIntegration_Conference_InvalidUUID(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.Conference(ctx, "00000000-0000-0000-0000-000000000000", "test-conf-room")
	if err == nil {
		t.Error("Expected error for invalid UUID conference")
	}
}

func TestIntegration_Eavesdrop_InvalidUUID(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.Eavesdrop(ctx, "00000000-0000-0000-0000-000000000000", "11111111-1111-1111-1111-111111111111")
	if err == nil {
		t.Error("Expected error for invalid UUID eavesdrop")
	}
}

func TestIntegration_Recording_InvalidUUID(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.StartRecording(ctx, "00000000-0000-0000-0000-000000000000", "/tmp/test.wav")
	if err == nil {
		t.Error("Expected error for invalid UUID recording")
	}
}

func TestIntegration_SofiaStatus(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// sofia status is used by trunk health monitor
	result, err := c.SendCommand(ctx, "sofia status")
	if err != nil {
		// mod_sofia may not be loaded in Docker containers without network access
		t.Skipf("sofia status not available (expected in Docker): %v", err)
	}
	t.Logf("Sofia status:\n%s", result)

	if !strings.Contains(result, "Name") {
		t.Error("Expected sofia status to contain profile information")
	}
}

func TestIntegration_SanitizeParam(t *testing.T) {
	// Unit test - doesn't need FreeSWITCH
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"with\nnewline", "withnewline"},
		{"with\rcarriage", "withcarriage"},
		{"with\x00null", "withnull"},
		{"1234#*", "1234#*"},
		{"user/1001@default", "user/1001@default"},
		{"+8613800138000", "+8613800138000"},
	}

	for _, tt := range tests {
		got := sanitizeParam(tt.input)
		if got != tt.expected {
			t.Errorf("sanitizeParam(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestIntegration_CircuitBreaker(t *testing.T) {
	// Unit test - doesn't need FreeSWITCH
	cb := &circuitBreaker{
		threshold:    3,
		state:        "closed",
		resetTimeout: 100 * time.Millisecond,
	}

	// Should allow in closed state
	if !cb.allow() {
		t.Error("Expected allow in closed state")
	}

	// Record failures up to threshold
	cb.recordFailure()
	cb.recordFailure()
	if !cb.allow() {
		t.Error("Expected allow before threshold")
	}

	cb.recordFailure()
	if cb.allow() {
		t.Error("Expected deny after threshold reached")
	}

	// Wait for reset timeout
	time.Sleep(150 * time.Millisecond)
	if !cb.allow() {
		t.Error("Expected allow after reset timeout (half-open)")
	}

	// Success should reset
	cb.recordSuccess()
	if cb.state != "closed" {
		t.Errorf("Expected closed state after success, got %s", cb.state)
	}
}

func TestIntegration_MultipleCommands(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Run multiple commands to test connection reuse
	for i := 0; i < 5; i++ {
		result, err := c.SendCommand(ctx, "status")
		if err != nil {
			t.Fatalf("Command %d failed: %v", i, err)
		}
		if !strings.Contains(result, "FreeSWITCH") {
			t.Errorf("Command %d: unexpected result: %s", i, result)
		}
	}
}

func TestIntegration_ShowChannels(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := c.SendCommand(ctx, "show channels")
	if err != nil {
		t.Fatalf("show channels failed: %v", err)
	}
	t.Logf("Active channels:\n%s", result)
}

func TestIntegration_GlobalGetvar(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test global_getvar used to check system variables
	result, err := c.SendCommand(ctx, "global_getvar hostname")
	if err != nil {
		t.Fatalf("global_getvar failed: %v", err)
	}
	t.Logf("Hostname: %s", result)
}

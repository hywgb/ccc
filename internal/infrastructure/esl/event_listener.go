package esl

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// EventHandler is invoked for each ESL event the listener receives.
// Implementations should be non-blocking; long work must run in its own goroutine.
type EventHandler interface {
	HandleESLEvent(ctx context.Context, ev Event)
}

// Event is a decoded ESL event with a small set of fields the lifecycle
// service cares about. Additional headers are kept in Headers for callers that
// need them.
type Event struct {
	Name        string
	ChannelUUID string
	Direction   string
	HangupCause string
	CallerID    string
	DestNumber  string
	Headers     map[string]string
	OccurredAt  time.Time
}

// EventListener subscribes to a configured set of ESL events on a dedicated
// long-lived connection and dispatches them to the registered handler.
type EventListener struct {
	cfg     Config
	logger  zerolog.Logger
	handler EventHandler
	events  []string
}

// NewEventListener constructs an ESL event listener. The handler is invoked
// for every event matching the events list (defaults to a minimal call-state set).
func NewEventListener(cfg Config, h EventHandler, events ...string) *EventListener {
	if len(events) == 0 {
		events = []string{
			"CHANNEL_ANSWER",
			"CHANNEL_HANGUP",
			"CHANNEL_HANGUP_COMPLETE",
			"CHANNEL_BRIDGE",
			"CHANNEL_PARK",
		}
	}
	return &EventListener{cfg: cfg, logger: cfg.Logger, handler: h, events: events}
}

// Run blocks until ctx is canceled, reconnecting with backoff on failure.
func (l *EventListener) Run(ctx context.Context) {
	backoff := time.Second
	for {
		if ctx.Err() != nil {
			return
		}
		if err := l.connectAndConsume(ctx); err != nil && !errors.Is(err, context.Canceled) {
			l.logger.Warn().Err(err).Dur("retry_in", backoff).Msg("ESL event listener disconnected")
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}
			if backoff < 30*time.Second {
				backoff *= 2
			}
			continue
		}
		backoff = time.Second
	}
}

func (l *EventListener) connectAndConsume(ctx context.Context) error {
	addr := net.JoinHostPort(l.cfg.Host, strconv.Itoa(l.cfg.Port))
	tcpConn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}
	defer tcpConn.Close()

	reader := bufio.NewReader(tcpConn)

	hdr, err := readHeaders(reader)
	if err != nil {
		return fmt.Errorf("read auth/request: %w", err)
	}
	if hdr["Content-Type"] != "auth/request" {
		return fmt.Errorf("expected auth/request, got %s", hdr["Content-Type"])
	}
	if _, err := fmt.Fprintf(tcpConn, "auth %s\n\n", l.cfg.Password); err != nil {
		return fmt.Errorf("send auth: %w", err)
	}
	hdr, err = readHeaders(reader)
	if err != nil {
		return fmt.Errorf("read auth reply: %w", err)
	}
	if !strings.HasPrefix(hdr["Reply-Text"], "+OK") {
		return fmt.Errorf("auth rejected: %s", hdr["Reply-Text"])
	}

	subscribe := "event plain " + strings.Join(l.events, " ") + "\n\n"
	if _, err := io.WriteString(tcpConn, subscribe); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	if _, err := readHeaders(reader); err != nil {
		return fmt.Errorf("read subscribe reply: %w", err)
	}
	l.logger.Info().Str("host", addr).Strs("events", l.events).Msg("ESL event listener connected")

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		frameHdr, err := readHeaders(reader)
		if err != nil {
			return fmt.Errorf("read frame headers: %w", err)
		}
		ctype := frameHdr["Content-Type"]
		clenStr := frameHdr["Content-Length"]
		if clenStr == "" {
			continue
		}
		clen, err := strconv.Atoi(clenStr)
		if err != nil {
			return fmt.Errorf("invalid content-length %q: %w", clenStr, err)
		}
		body := make([]byte, clen)
		if _, err := io.ReadFull(reader, body); err != nil {
			return fmt.Errorf("read body: %w", err)
		}
		if ctype != "text/event-plain" {
			continue
		}
		ev := parseEvent(body)
		if ev.Name == "" {
			continue
		}
		l.handler.HandleESLEvent(ctx, ev)
	}
}

func parseEvent(body []byte) Event {
	ev := Event{Headers: make(map[string]string), OccurredAt: time.Now()}
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		idx := strings.Index(line, ": ")
		if idx <= 0 {
			continue
		}
		k := line[:idx]
		v := line[idx+2:]
		if decoded, err := url.QueryUnescape(v); err == nil {
			v = decoded
		}
		ev.Headers[k] = v
	}
	ev.Name = ev.Headers["Event-Name"]
	ev.ChannelUUID = ev.Headers["Unique-ID"]
	if ev.ChannelUUID == "" {
		ev.ChannelUUID = ev.Headers["Channel-Call-UUID"]
	}
	ev.Direction = ev.Headers["Call-Direction"]
	ev.HangupCause = ev.Headers["Hangup-Cause"]
	ev.CallerID = ev.Headers["Caller-Caller-ID-Number"]
	ev.DestNumber = ev.Headers["Caller-Destination-Number"]
	return ev
}

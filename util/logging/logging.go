// Package logging provides a colorized slog handler for local development.
//
// Usage:
//
//	logging.Setup() // installs as slog.Default()
//
// The LOG_FORMAT environment variable controls the output format:
//   - "json": structured JSON output (production)
//   - default: colorized text output (development)
package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

// ANSI color codes for log levels (background colors).
const (
	colorReset = "\033[0m"
	colorGray  = "\033[90m"

	// Background colors with white text for level badges.
	bgRed    = "\033[41;97m"  // red bg, white text
	bgYellow = "\033[43;30m"  // yellow bg, black text
	bgGreen  = "\033[42;30m"  // green bg, black text
	bgGray   = "\033[100;97m" // gray bg, white text
)

// Setup initializes the global slog logger.
//
// If LOG_FORMAT=json, it uses slog.JSONHandler.
// Otherwise, it uses a colorized text handler for TTY output.
func Setup() {
	var handler slog.Handler

	if os.Getenv("LOG_FORMAT") == "json" {
		handler = slog.NewJSONHandler(os.Stderr, nil)
	} else {
		handler = NewColorHandler(os.Stderr)
	}

	slog.SetDefault(slog.New(handler))
}

// ColorHandler outputs logs with a colored level badge followed by gray text.
type ColorHandler struct {
	writer io.Writer
	attrs  []slog.Attr
	group  string
	mu     sync.Mutex
}

// NewColorHandler creates a handler that outputs colorized text logs.
func NewColorHandler(w io.Writer) *ColorHandler {
	return &ColorHandler{
		writer: w,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *ColorHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// Handle formats and writes the log record with colored level badge.
func (h *ColorHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Format: [LEVEL] time=... msg=... key=value...
	levelText := r.Level.String()
	bg := levelBgColor(r.Level)

	// Write colored level badge.
	fmt.Fprintf(h.writer, "%s %s %s ", bg, levelText, colorReset)

	// Write the rest in gray.
	fmt.Fprintf(h.writer, "%stime=%s msg=%q", colorGray, r.Time.Format(time.RFC3339), r.Message)

	// Write pre-defined attrs.
	for _, a := range h.attrs {
		h.writeAttr(a)
	}

	// Write record attrs.
	r.Attrs(func(a slog.Attr) bool {
		h.writeAttr(a)
		return true
	})

	fmt.Fprintf(h.writer, "%s\n", colorReset)
	return nil
}

func (h *ColorHandler) writeAttr(a slog.Attr) {
	if a.Key == "" {
		return
	}
	key := a.Key
	if h.group != "" {
		key = h.group + "." + key
	}
	fmt.Fprintf(h.writer, " %s=%v", key, a.Value.Any())
}

// WithAttrs returns a new handler with the given attributes.
func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs), len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	newAttrs = append(newAttrs, attrs...)
	return &ColorHandler{
		writer: h.writer,
		attrs:  newAttrs,
		group:  h.group,
	}
}

// WithGroup returns a new handler with the given group name.
func (h *ColorHandler) WithGroup(name string) slog.Handler {
	newGroup := name
	if h.group != "" {
		newGroup = h.group + "." + name
	}
	return &ColorHandler{
		writer: h.writer,
		attrs:  h.attrs,
		group:  newGroup,
	}
}

// levelBgColor returns the ANSI background color code for the given log level.
func levelBgColor(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return bgRed
	case level >= slog.LevelWarn:
		return bgYellow
	case level >= slog.LevelInfo:
		return bgGreen
	default:
		return bgGray
	}
}

package sslog

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/slog"
)

// PipelineHandler supports handling a record through multiple handlers
type PipelineHandler []slog.Handler

var PipelineHalted = fmt.Errorf("pipeline handler halted")

var _ slog.Handler = PipelineHandler{}

// Enabled reports false only if all handlers are disabled or if no handlers are present
func (handlers PipelineHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	for _, handler := range handlers {
		if handler.Enabled(ctx, lvl) {
			return true
		}
	}

	return false
}

// Handle handles the Record.
// It will call Handle for each handler in the list of handlers.
// The Handle method specific to a handler will only be called if the corresponding Enabled returns true.
// If a handler in the list of handlers returns an error, the execution of subsequent handlers is skipped and an error is returned, however if the error returned is PipelineHalted, no error will be returned by PipelineHandler Handle method.
func (handlers PipelineHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, handler := range handlers {
		if !handler.Enabled(ctx, record.Level) {
			continue
		}

		err := handler.Handle(ctx, record)
		if err != nil && !errors.Is(err, PipelineHalted) {
			return fmt.Errorf("pipeline interrupted: %w", err)
		}
		if err != nil && errors.Is(err, PipelineHalted) {
			return nil
		}
	}

	return nil
}

// WithAttrs is applied to all nested handlers and returns a new copy of BroadcastHandler
func (handlers PipelineHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make(PipelineHandler, len(handlers))

	for i, handler := range handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}

	return newHandlers
}

// WithGroup is applied to all nested handlers and returns a new copy of BroadcastHandler
// If the name is empty, WithGroup returns the receiver.
func (handlers PipelineHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return handlers
	}

	newHandlers := make(PipelineHandler, len(handlers))

	for i, handler := range handlers {
		newHandlers[i] = handler.WithGroup(name)
	}

	return newHandlers
}

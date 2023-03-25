package sslog

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/slog"
)

var ErrPipelineInterrupted = fmt.Errorf("pipline interrupted")
var ErrPipelineHalted = fmt.Errorf("pipeline halted")

type ErrPipelineContinue struct {
	Record slog.Record
}

func (ErrPipelineContinue) Error() string {
	return "pipeline must continue with attached record"
}

func PipelineContinue(record slog.Record) ErrPipelineContinue {
	return ErrPipelineContinue{record}
}

// PipelineHandler supports handling a record through multiple handlers
type PipelineHandler []slog.Handler

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
	currentRecord := record

	for _, handler := range handlers {
		// Skip disabled handler
		if !handler.Enabled(ctx, currentRecord.Level) {
			continue
		}

		err := handler.Handle(ctx, currentRecord)
		var errPipelineContinue ErrPipelineContinue
		// Halt, which means interrupted voluntarily, not an exception
		if err != nil && errors.Is(err, ErrPipelineHalted) {
			return nil
		}
		// Continue, which means continue pipeline with modified record
		if err != nil && errors.As(err, &errPipelineContinue) {
			currentRecord = errPipelineContinue.Record
			continue
		}
		// Interrupt, which means **exceptional**
		if err != nil {
			return fmt.Errorf("%w %w", ErrPipelineInterrupted, err)
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

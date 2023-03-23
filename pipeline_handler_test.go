package sslog_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/thehungry-dev/sslog"
	"golang.org/x/exp/slog"
)

func TestPipelineHandlerEnabled(t *testing.T) {
	t.Parallel()

	t.Run("false when empty", func(t *testing.T) {
		pipeline := sslog.PipelineHandler{}

		enabled := pipeline.Enabled(context.Background(), slog.LevelDebug)

		if enabled {
			t.Error()
		}
	})

	t.Run("false when all handlers Enabled return false", func(t *testing.T) {
		mockHandler1 := &sslog.MockHandler{Active: false}
		mockHandler2 := &sslog.MockHandler{Active: false}
		pipeline := sslog.PipelineHandler{mockHandler1, mockHandler2}

		enabled := pipeline.Enabled(context.Background(), slog.LevelDebug)

		if enabled {
			t.Error()
		}
	})

	t.Run("true when at least one handler is Enabled", func(t *testing.T) {
		mockHandler1 := &sslog.MockHandler{Active: false}
		mockHandler2 := &sslog.MockHandler{Active: true}
		pipeline := sslog.PipelineHandler{mockHandler1, mockHandler2}

		enabled := pipeline.Enabled(context.Background(), slog.LevelDebug)

		if !enabled {
			t.Error()
		}
	})
}

func TestPipelineHandlerHandle(t *testing.T) {
	t.Parallel()

	t.Run("returns error when at least one handler returns error", func(t *testing.T) {
		mockErr := fmt.Errorf("mock error")
		mockHandler := sslog.ActiveMockHandler()
		mockHandlerErr := sslog.MockHandler{HandleError: mockErr, Active: true}
		pipeline := sslog.PipelineHandler{mockHandler, &mockHandlerErr}

		err := pipeline.Handle(context.Background(), slog.Record{})

		if err == nil {
			t.Error("error must be returned from Handle")
		}
		if !errors.Is(err, mockErr) {
			t.Errorf("invalid mock error: %s", err.Error())
		}
	})

	t.Run("skips following handlers when at least one returns error", func(t *testing.T) {
		mockErr := fmt.Errorf("mock error")
		mockHandlerErr := sslog.MockHandler{HandleError: mockErr, Active: true}
		mockHandler := sslog.ActiveMockHandler()
		pipeline := sslog.PipelineHandler{&mockHandlerErr, mockHandler}

		err := pipeline.Handle(context.Background(), slog.Record{})

		if err == nil {
			t.Error("error must be returned from Handle")
		}
		if !errors.Is(err, mockErr) {
			t.Errorf("invalid mock error: %s", err.Error())
		}
		if len(mockHandler.Recordings) > 0 {
			t.Error("second handler Handle was called when it should have been interrupted")
		}
	})

	t.Run("skips following handlers when at least one is halted without returning error", func(t *testing.T) {
		mockHandlerErr := sslog.MockHandler{HandleError: sslog.PipelineHalted, Active: true}
		mockHandler := sslog.ActiveMockHandler()
		pipeline := sslog.PipelineHandler{&mockHandlerErr, mockHandler}

		err := pipeline.Handle(context.Background(), slog.Record{})

		if err != nil {
			t.Errorf("halted pipeline must have no error: %s", err.Error())
		}
		if len(mockHandler.Recordings) > 0 {
			t.Error("second handler Handle was called when it should have been interrupted")
		}
	})

	t.Run("calls all handlers Handle when no error is returned", func(t *testing.T) {
		mockHandler1 := sslog.ActiveMockHandler()
		mockHandler2 := sslog.ActiveMockHandler()
		pipeline := sslog.PipelineHandler{mockHandler1, mockHandler2}

		err := pipeline.Handle(context.Background(), slog.Record{})

		if err != nil {
			t.Errorf("handle must return no errors: %s", err.Error())
		}
		if len(mockHandler1.Recordings) != 1 {
			t.Errorf("each handler must be invoked once: %+v", mockHandler1.Recordings)
		}
		if len(mockHandler2.Recordings) != 1 {
			t.Errorf("each handler must be invoked once: %+v", mockHandler2.Recordings)
		}
	})
}

package handlers

import (
	"context"

	"golang.org/x/exp/slog"
)

const (
	TrackHandle = iota
	TrackWithAttrs
	TrackWithGroup
)

type Recording struct {
	RecordingType int
	Record        slog.Record
	Attrs         []slog.Attr
	Name          string
}

// MockHandler records every call to the handler.
// HandleError can be set to force Handle to return this error.
type MockHandler struct {
	Recordings  []Recording
	Active      bool
	HandleError error
}

var _ slog.Handler = &MockHandler{}

func ActiveMockHandler() *MockHandler {
	return &MockHandler{Active: true}
}

func (handler *MockHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return handler.Active
}

// Handle halts the pipeline if the record with name "sslogTags" attribute is rejected
func (handler *MockHandler) Handle(_ context.Context, record slog.Record) error {
	handler.Recordings = append(handler.Recordings, Recording{RecordingType: TrackHandle, Record: record})

	return handler.HandleError
}

// WithAttrs returns the receiver
func (handler *MockHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handler.Recordings = append(handler.Recordings, Recording{RecordingType: TrackWithAttrs, Attrs: attrs})

	return handler
}

// WithGroup returns the receiver
func (handler *MockHandler) WithGroup(name string) slog.Handler {
	handler.Recordings = append(handler.Recordings, Recording{RecordingType: TrackWithGroup, Name: name})

	return handler
}

// Clear empties recordings
func (handler *MockHandler) Clear() {
	handler.Recordings = []Recording{}
}

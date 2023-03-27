package handlers_test

import (
	"context"
	"testing"

	"github.com/thehungry-dev/sslog/pkg/ctrls"
	"github.com/thehungry-dev/sslog/pkg/handlers"
)

func TestTagFilterHandlerHandle(t *testing.T) {
	t.Parallel()

	t.Run("continues pipeline when filter matches record tags", func(t *testing.T) {
		record := ctrls.RecordExample()
		mockHandler := handlers.ActiveMockHandler()
		tagFilterHandler := ctrls.TagFilterHandlerMatchingExample()
		pipeline := handlers.PipelineHandler{tagFilterHandler, mockHandler}
		t.Logf("tag filter handler configuration: %s", tagFilterHandler.String())

		pipeline.Handle(context.Background(), record)

		if len(mockHandler.Recordings) == 0 {
			t.Errorf("record has been filtered out: %+v", record)
		}
		if mockHandler.Recordings[0].RecordingType != handlers.TrackHandle {
			t.Errorf("unknown recording: %+v", mockHandler.Recordings[0])
		}
		if mockHandler.Recordings[0].Record.Message != record.Message {
			t.Logf("expected record: %+v", record)
			t.Errorf("unknown record tracked: %+v", mockHandler.Recordings[0].Record)
		}
	})

	t.Run("halts pipeline when filter rejects record tags", func(t *testing.T) {
		record := ctrls.RecordExample()
		mockHandler := handlers.ActiveMockHandler()
		tagFilterHandler := ctrls.TagFilterHandlerRejectingExample()
		pipeline := handlers.PipelineHandler{tagFilterHandler, mockHandler}
		t.Logf("tag filter handler configuration: %s", tagFilterHandler.String())

		pipeline.Handle(context.Background(), record)

		if len(mockHandler.Recordings) != 0 {
			t.Errorf("record should have been filtered out: %+v", mockHandler.Recordings[0])
		}
	})
}

package main

import (
	"os"

	"github.com/thehungry-dev/sslog"
	"github.com/thehungry-dev/sslog/pkg/handlers"
	"golang.org/x/exp/slog"
)

func main() {
	textHandler := slog.NewTextHandler(os.Stderr)
	filterTagHandler := handlers.ParseTagFilter("+requiredTag,-excludedTag")
	pipeline := handlers.PipelineHandler{filterTagHandler, textHandler}
	logger := slog.
		New(pipeline).
		With(sslog.Tags("anotherTag", "requiredTag")).
		With("foo", "bar").
		With("foo", "whatever").
		With(sslog.Tags("whateverTag"))
	slog.SetDefault(logger)

	slog.Info("message 1")
}

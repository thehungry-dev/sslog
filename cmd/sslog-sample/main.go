package main

import (
	"os"

	"github.com/thehungry-dev/sslog"
	"golang.org/x/exp/slog"
)

func main() {
	textHandler := slog.NewTextHandler(os.Stderr)
	filterTagHandler := sslog.ParseTagFilter("+requiredTag,-excludedTag")
	pipeline := sslog.PipelineHandler{filterTagHandler, textHandler}
	logger := slog.
		New(pipeline).
		With(sslog.Tags("anotherTag", "requiredTag")).
		With("foo", "bar").
		With(sslog.Tags("whateverTag"))
	slog.SetDefault(logger)

	slog.Info("message 1")
}

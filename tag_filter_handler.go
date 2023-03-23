package sslog

import (
	"context"
	"fmt"

	"github.com/thehungry-dev/rag"
	"golang.org/x/exp/slog"
)

var TagsKey string = "sslogTags"

// Filter out records that for which the TagFilter doesn't accept the tags.
// If no tags are attached to the record, the record is treated as with an empty slice of tags.
type TagFilterHandler struct {
	tagFilter *rag.TagFilter
	Active    bool
}

var _ slog.Handler = TagFilterHandler{}

// ParseTagFilter initializes an Active TagFilterHandler with the provided *rag.TagFilter string
func ParseTagFilter(filterText string) TagFilterHandler {
	tagFilter := rag.Parse(filterText)

	return TagFilterHandler{tagFilter: tagFilter, Active: true}
}

// Enabled returns false when Active is false, otherwise it returns true
func (handler TagFilterHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return handler.Active
}

// Handle halts the pipeline if the record with name "sslogTags" attribute is rejected
func (handler TagFilterHandler) Handle(_ context.Context, record slog.Record) error {
	var foundTags bool
	var untypedTags interface{} = []string{}

	record.Attrs(func(attr slog.Attr) {
		if foundTags {
			return
		}
		// Attribute is not the list of tags
		if attr.Key != TagsKey {
			return
		}

		foundTags = true
		untypedTags = attr.Value.Any()
	})

	tags, ok := untypedTags.([]string)
	if !ok {
		return fmt.Errorf("record attribute value for tags is not a slice of strings")
	}

	// Ignore record if tags are not accepted
	if handler.tagFilter.Reject(tags) {
		return PipelineHalted
	}

	return nil
}

// WithAttrs returns the receiver
func (handler TagFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return handler
}

// WithGroup returns the receiver
func (handler TagFilterHandler) WithGroup(name string) slog.Handler {
	return handler
}

func (handler TagFilterHandler) String() string {
	return fmt.Sprintf("{Enabled: %t, FilterText: \"%s\"}", handler.Active, handler.tagFilter.String())
}

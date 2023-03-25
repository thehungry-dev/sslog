package sslog

import (
	"context"
	"fmt"

	"github.com/thehungry-dev/rag"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

var TagsKey string = "tags"

// Filter out records that for which the TagFilter doesn't accept the tags.
// If no tags are attached to the record, the record is treated as with an empty slice of tags.
type TagFilterHandler struct {
	tagFilter   *rag.TagFilter
	Active      bool
	defaultTags []string
}

var _ slog.Handler = TagFilterHandler{}

func addAttributesToTags(tags []string, attrs ...slog.Attr) ([]string, uint) {
	var tagsAdded uint

	for _, attr := range attrs {
		// Attribute is not the list of tags
		if attr.Key != TagsKey {
			continue
		}

		// Attribute is of incorrect type
		additionalTags, ok := attr.Value.Any().([]string)
		if !ok {
			continue
		}

		tagsAdded = tagsAdded + uint(len(additionalTags))
		tags = append(tags, additionalTags...)
	}

	return tags, tagsAdded
}

func getRecordTags(record slog.Record) ([]string, uint) {
	var totalTags uint
	recordTags := []string{}

	record.Attrs(func(attr slog.Attr) {
		var addedTags uint
		recordTags, addedTags = addAttributesToTags(recordTags, attr)
		totalTags = totalTags + addedTags
	})

	return recordTags, totalTags
}

func getRecordTagsWithDefault(record slog.Record, defaultTags []string) []string {
	newTags := slices.Clone(defaultTags)
	recordTags, _ := getRecordTags(record)

	newTags = append(newTags, recordTags...)

	return newTags
}

func resetRecordTags(record slog.Record, tags []string) slog.Record {
	newRecord := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)

	record.Attrs(func(attr slog.Attr) {
		// Skip tags attribute
		if attr.Key == TagsKey {
			return
		}

		record.AddAttrs(attr)
	})

	// Add modified tag attribute
	newRecord.Add(TagsKey, tags)

	return newRecord
}

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
	tags := getRecordTagsWithDefault(record, handler.defaultTags)

	// Ignore record if tags are not accepted
	rejected := handler.tagFilter.Reject(tags)
	if rejected {
		return ErrPipelineHalted
	}

	newRecord := resetRecordTags(record, tags)

	return PipelineContinue(newRecord)
}

// WithAttrs returns a copy of the handler only if the attribute key follows sslog.TagsKey.
// If it is sslog.TagsKey, the attribute must be a slice of strings which are appended to each record's tags.
// Otherwise, returns the receiver
func (handler TagFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	defaultTags := slices.Clone(handler.defaultTags)

	defaultTags, tagsAdded := addAttributesToTags(defaultTags, attrs...)
	if tagsAdded == 0 {
		return handler
	}

	newHandler := TagFilterHandler{
		tagFilter:   handler.tagFilter,
		Active:      handler.Active,
		defaultTags: defaultTags,
	}

	return newHandler
}

// WithGroup returns the receiver
func (handler TagFilterHandler) WithGroup(name string) slog.Handler {
	return handler
}

func (handler TagFilterHandler) String() string {
	return fmt.Sprintf("{Enabled: %t, FilterText: \"%s\"}", handler.Active, handler.tagFilter.String())
}

// Package ctrls helps generating sample data
package ctrls

import (
	"time"

	ragctrls "github.com/thehungry-dev/rag/pkg/ctrls"
	"github.com/thehungry-dev/sslog"
	"golang.org/x/exp/slog"
)

func TagsExample() []string              { return ragctrls.TagsNonMatchingMissingAllRequiredOneOfExample() }
func FilterTextMatchingExample() string  { return ragctrls.StringExcludedRequiredExample() }
func FilterTextRejectingExample() string { return ragctrls.StringExample() }
func TagFilterHandlerMatchingExample() sslog.TagFilterHandler {
	return sslog.ParseTagFilter(FilterTextMatchingExample())
}
func TagFilterHandlerRejectingExample() sslog.TagFilterHandler {
	return sslog.ParseTagFilter(FilterTextRejectingExample())
}

func RecordDate() time.Time { return time.Date(2000, time.January, 1, 1, 1, 1, 1, time.UTC) }
func RecordExample() slog.Record {
	record := slog.NewRecord(RecordDate(), slog.LevelInfo, "a message", 0)

	record.Add(sslog.TagsKey, TagsExample())

	return record
}

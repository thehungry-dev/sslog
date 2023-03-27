// Package sslog provides a backend for slog that allows filtering out log statements based on tags
package sslog

var TagsKey string = "tags"

func Tags(tags ...string) (string, []string) {
	return TagsKey, tags
}

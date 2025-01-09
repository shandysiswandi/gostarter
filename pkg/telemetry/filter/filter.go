package filter

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/shandysiswandi/gostarter/pkg/codec"
)

type Filter struct {
	headers []string
	queries []string
	fields  []string
	json    codec.Codec
}

// NewFilter creates a new Filter instance with the provided options.
// It initializes the JSON codec and sets up empty slices for headers, queries, and fields.
func NewFilter(opts ...OptionFilter) *Filter {
	f := &Filter{
		json:    codec.NewJSONCodec(),
		headers: make([]string, 0),
		queries: make([]string, 0),
		fields:  make([]string, 0),
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

// Query takes a raw URL string and replaces the values of specified query parameters
// with asterisks (***). Special characters in the query parameter keys are escaped to ensure
// proper matching.
func (f *Filter) Query(rawURL string) string {
	escapedKeys := make([]string, 0, len(f.queries))
	for _, key := range f.queries {
		escapedKeys = append(escapedKeys, regexp.QuoteMeta(key)) // Escape any special characters in keys.
	}
	pattern := fmt.Sprintf(`(?i)(%s)=.*?(&|$)`, strings.Join(escapedKeys, "|"))

	return regexp.MustCompile(pattern).ReplaceAllString(rawURL, `$1=***$2`)
}

// Body processes the given JSON body and returns a map representation of it.
// If the body is empty or if there is an error during decoding, an empty map is returned.
func (f *Filter) Body(body []byte) map[string]any {
	if len(body) == 0 {
		return map[string]any{}
	}

	var tBody map[string]any
	if err := f.json.Decode(body, &tBody); err != nil {
		return map[string]any{}
	}

	f.processBody(tBody)

	return tBody
}

// processBody masks sensitive fields in the provided data map.
// It recursively processes nested maps and slices to ensure all sensitive fields are masked.
func (f *Filter) processBody(data map[string]any) {
	mask := "***"
	for key, value := range data {
		if slices.Contains(f.fields, key) {
			data[key] = mask
		} else if nestedMap, ok := value.(map[string]any); ok {
			// If the value is a nested map, recursively mask it
			f.processBody(nestedMap)
		} else if nestedSlice, ok := value.([]any); ok {
			// If the value is a slice, iterate and handle nested maps
			for i, item := range nestedSlice {
				if nestedItem, ok := item.(map[string]any); ok {
					f.processBody(nestedItem)
					nestedSlice[i] = nestedItem
				}
			}
		}
	}
}

// Header processes the provided HTTP headers and returns a map with the headers
// in lowercase and their values joined by commas. If a header key matches any
// of the filters in the Filter's headers, its value is replaced with "***".
func (f *Filter) Header(hh map[string][]string) map[string]string {
	if len(hh) == 0 {
		return map[string]string{}
	}

	meta := make(map[string]string)
	for key, value := range hh {
		meta[strings.ToLower(key)] = strings.Join(value, ",")
	}

	for _, filter := range f.headers {
		if _, ok := meta[filter]; ok {
			meta[filter] = "***"
		}
	}

	return meta
}

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

func (f *Filter) Query(rawURL string) string {
	escapedKeys := make([]string, len(f.queries))
	for i, key := range f.queries {
		escapedKeys[i] = regexp.QuoteMeta(key) // Escape any special characters in keys.
	}
	pattern := fmt.Sprintf(`(?i)(%s)=.*?(&|$)`, strings.Join(escapedKeys, "|"))

	return regexp.MustCompile(pattern).ReplaceAllString(rawURL, `$1=***$2`)
}

func (f *Filter) Body(body []byte) map[string]any {
	if body == nil {
		return map[string]any{}
	}

	var tBody map[string]any
	if err := f.json.Decode(body, &tBody); err != nil {
		return map[string]any{}
	}

	f.processBody(tBody)

	return tBody
}

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

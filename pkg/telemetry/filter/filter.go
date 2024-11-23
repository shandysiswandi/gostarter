package filter

import (
	"strings"

	"github.com/shandysiswandi/gostarter/pkg/codec"
	"google.golang.org/grpc/metadata"
)

type Filter struct {
	headers []string
	json    codec.Codec
}

type OptionFilter func(*Filter)

func WithHeaders(header ...string) OptionFilter {
	return func(f *Filter) {
		f.headers = append(f.headers, header...)
	}
}

func NewFilter(opts ...OptionFilter) *Filter {
	f := &Filter{
		json:    codec.NewJSONCodec(),
		headers: make([]string, 0),
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func (f *Filter) Metadata(md metadata.MD) string {
	if len(md) == 0 {
		return ""
	}

	meta := make(map[string]string)
	for key, value := range md {
		meta[key] = strings.Join(value, ",")
	}

	for _, filter := range f.headers {
		if _, ok := meta[filter]; ok {
			meta[filter] = "***"
		}
	}

	bt, err := f.json.Encode(meta)
	if err != nil {
		return ""
	}

	return string(bt)
}

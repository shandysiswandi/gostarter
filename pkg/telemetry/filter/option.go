package filter

type OptionFilter func(*Filter)

func WithHeaders(header ...string) OptionFilter {
	return func(f *Filter) {
		f.headers = append(f.headers, header...)
	}
}

func WithQueries(query ...string) OptionFilter {
	return func(f *Filter) {
		f.queries = append(f.queries, query...)
	}
}

func WithFields(field ...string) OptionFilter {
	return func(f *Filter) {
		f.fields = append(f.fields, field...)
	}
}

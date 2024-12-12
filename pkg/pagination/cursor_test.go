package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCursorBased(t *testing.T) {
	type args struct {
		cursor string
		limit  string
	}
	tests := []struct {
		name       string
		args       args
		wantCursor uint64
		wantLimit  int
	}{
		{
			name: "ErrorParseLimitAndCursorEmpty",
			args: args{
				cursor: "",
				limit:  "Atoi",
			},
			wantCursor: 0,
			wantLimit:  DefaultLimit,
		},
		{
			name: "LimitLessThanOneAndCursorErrDecode",
			args: args{
				cursor: "-",
				limit:  "0",
			},
			wantCursor: 0,
			wantLimit:  DefaultLimit,
		},
		{
			name: "LimitMoreThanMaxAndCursorErrParse",
			args: args{
				cursor: "atoi",
				limit:  "101",
			},
			wantCursor: 0,
			wantLimit:  MaxLimit,
		},
		{
			name: "Success",
			args: args{
				cursor: "NTY",
				limit:  "44",
			},
			wantCursor: 56,
			wantLimit:  44,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cursor, limit := ParseCursorBased(tt.args.cursor, tt.args.limit)
			assert.Equal(t, tt.wantCursor, cursor)
			assert.Equal(t, tt.wantLimit, limit)
		})
	}
}

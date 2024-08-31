package shortly

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		dep     Dependency
		want    *Expose
		wantErr error
	}{
		{
			name: "Success",
			dep: Dependency{
				Router: &httprouter.Router{},
			},
			want:    &Expose{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := New(tt.dep)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

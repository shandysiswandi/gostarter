package region

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	cMock "github.com/shandysiswandi/gostarter/pkg/config/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		dep     func() Dependency
		want    *Expose
		wantErr error
	}{
		{
			name: "Success",
			dep: func() Dependency {
				mc := cMock.NewMockConfig(t)

				mc.EXPECT().GetString("database.driver").Return("mysql").Once()

				return Dependency{
					Config: mc,
					Router: &httprouter.Router{},
				}
			},
			want:    &Expose{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := New(tt.dep())
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

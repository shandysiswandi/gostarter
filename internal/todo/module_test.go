package todo

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/job"
	configMock "github.com/shandysiswandi/gostarter/pkg/config/mocker"
	"github.com/shandysiswandi/gostarter/pkg/task"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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
				mc := new(configMock.MockConfig)

				mc.EXPECT().GetString("database.driver").Return("mysql").Once()
				mc.EXPECT().GetBool("feature.flag.graphql.playground").Return(true).Once()

				return Dependency{Config: mc, Router: &httprouter.Router{}, GRPCServer: grpc.NewServer()}
			},
			want:    &Expose{Tasks: []task.Runner{&job.ExampleJob{}}},
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

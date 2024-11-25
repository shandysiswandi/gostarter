package todo

import (
	"testing"

	configMock "github.com/shandysiswandi/gostarter/pkg/config/mocker"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		dep     func() Dependency
		wantErr error
	}{
		{
			name: "Success",
			dep: func() Dependency {
				mc := new(configMock.MockConfig)

				mc.EXPECT().GetString("database.driver").Return("mysql").Once()
				mc.EXPECT().GetBool("feature.flag.graphql.playground").Return(true).Once()
				mc.EXPECT().GetBool("feature.flag.todo.job").Return(true).Once()

				return Dependency{
					Config:     mc,
					Router:     framework.NewRouter(),
					GQLRouter:  framework.NewRouter(),
					GRPCServer: grpc.NewServer(),
				}
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := New(tt.dep())
			assert.Equal(t, tt.wantErr, err)
			assert.NotNil(t, got)
		})
	}
}

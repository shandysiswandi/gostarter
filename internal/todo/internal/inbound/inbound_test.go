package inbound

import (
	"testing"

	mockConfig "github.com/shandysiswandi/gostarter/pkg/config/mocker"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"google.golang.org/grpc"
)

func TestInbound_RegisterTodoServiceServer(t *testing.T) {
	tests := []struct {
		name string
		in   func() Inbound
	}{
		{
			name: "Success",
			in: func() Inbound {
				configMock := mockConfig.NewMockConfig(t)

				configMock.EXPECT().
					GetBool("feature.flag.graphql.playground").
					Return(true)

				return Inbound{
					Config:         configMock,
					Router:         framework.New(),
					GQLRouter:      framework.New(),
					GRPCServer:     grpc.NewServer(),
					CodecJSON:      nil,
					CreateUC:       nil,
					DeleteUC:       nil,
					FindUC:         nil,
					FetchUC:        nil,
					UpdateStatusUC: nil,
					UpdateUC:       nil,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.in().RegisterTodoServiceServer()
		})
	}
}

package inbound

import (
	"testing"

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

				return Inbound{
					Router:         framework.NewRouter(),
					GQLRouter:      framework.NewRouter(),
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

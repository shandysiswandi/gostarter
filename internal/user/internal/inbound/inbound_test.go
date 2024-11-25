package inbound

import (
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/framework"
	"google.golang.org/grpc"
)

func TestInbound_RegisterUserServiceServer(t *testing.T) {
	tests := []struct {
		name string
		in   Inbound
	}{
		{
			name: "Success",
			in: Inbound{
				Router:     framework.NewRouter(),
				GRPCServer: grpc.NewServer(),
				ProfileUC:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.in.RegisterUserServiceServer()
		})
	}
}

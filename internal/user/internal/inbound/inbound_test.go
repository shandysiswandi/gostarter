package inbound

import (
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/framework"
)

func TestInbound_RegisterUserServiceServer(t *testing.T) {
	tests := []struct {
		name string
		in   Inbound
	}{
		{
			name: "Success",
			in: Inbound{
				Router:    framework.NewRouter(),
				ProfileUC: nil,
				UpdateUC:  nil,
				LogoutUC:  nil,
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

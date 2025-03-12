package inbound

import (
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/framework"
)

func TestInbound_RegisterAuthServiceServer(t *testing.T) {
	tests := []struct {
		name string
		in   Inbound
	}{
		{
			name: "Success",
			in: Inbound{
				Router:           framework.NewRouter(),
				LoginUC:          nil,
				RegisterUC:       nil,
				RefreshTokenUC:   nil,
				ForgotPasswordUC: nil,
				ResetPasswordUC:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.in.RegisterAuthServiceServer()
		})
	}
}

package inbound

import (
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/framework"
)

func TestInbound_RegisterPaymentServiceServer(t *testing.T) {
	tests := []struct {
		name string
		in   func() Inbound
	}{
		{
			name: "Success",
			in: func() Inbound {
				return Inbound{
					Router: framework.NewRouter(),
					//
					PaymentTopupUC: nil,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.in().RegisterPaymentServiceServer()
		})
	}
}

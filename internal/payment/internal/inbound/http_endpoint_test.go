package inbound

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_httpEndpoint_PaymentTopup(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			c: func() framework.Context {
				body := bytes.NewBufferString("fake request")
				c := framework.NewTestContext(http.MethodPost, "/payments/topup", body)

				return c.Build()
			},
			want:    nil,
			wantErr: errInvalidBody,
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "payment.inbound.httpEndpoint.PaymentTopup")
				defer span.End()

				return &httpEndpoint{
					tel: tel,
				}
			},
		},
		{
			name: "ErrorParseAmount",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"reference_id":"uuid", "amount":"zzz"}`)
				c := framework.NewTestContext(http.MethodPost, "/payments/topup", body)

				return c.Build()
			},
			want:    nil,
			wantErr: errInvalidBody,
			mockFn: func(ctx context.Context) *httpEndpoint {
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "payment.inbound.httpEndpoint.PaymentTopup")
				defer span.End()

				return &httpEndpoint{
					tel: tel,
				}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"reference_id":"uuid", "amount":"1000.00"}`)
				c := framework.NewTestContext(http.MethodPost, "/payments/topup", body)

				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				ptMock := mockz.NewMockPaymentTopup(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "payment.inbound.httpEndpoint.PaymentTopup")
				defer span.End()

				amo, _ := decimal.NewFromString("1000.00")
				in := domain.PaymentTopupInput{
					ReferenceID: "uuid",
					Amount:      amo,
				}
				ptMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					tel:            tel,
					paymentTopupUC: ptMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"reference_id":"uuid", "amount":"1000.00"}`)
				c := framework.NewTestContext(http.MethodPost, "/payments/topup", body)

				return c.Build()
			},
			want: PaymentTopupResponse{
				ReferenceID: "uuid",
				Amount:      "1000.00",
				Balance:     "2000.00",
			},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				ptMock := mockz.NewMockPaymentTopup(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "payment.inbound.httpEndpoint.PaymentTopup")
				defer span.End()

				amo, _ := decimal.NewFromString("1000.00")
				in := domain.PaymentTopupInput{
					ReferenceID: "uuid",
					Amount:      amo,
				}
				out := &domain.PaymentTopupOutput{
					ReferenceID: in.ReferenceID,
					Amount:      in.Amount,
					Balance:     in.Amount.Add(in.Amount),
				}
				ptMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					tel:            tel,
					paymentTopupUC: ptMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.PaymentTopup(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

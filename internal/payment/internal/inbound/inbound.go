package inbound

import (
	"net/http"

	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
)

type Inbound struct {
	Router    *framework.Router
	Telemetry *telemetry.Telemetry
	//
	PaymentTopupUC domain.PaymentTopup
}

func (in Inbound) RegisterPaymentServiceServer() {
	he := &httpEndpoint{
		tel: in.Telemetry,
		//
		paymentTopupUC: in.PaymentTopupUC,
	}

	in.Router.Endpoint(http.MethodPost, "/payments/topup", he.PaymentTopup)
}

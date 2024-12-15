package inbound

import (
	"net/http"

	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
)

type Inbound struct {
	Router *framework.Router
	//
	PaymentTopupUC domain.PaymentTopup
}

func (in Inbound) RegisterPaymentServiceServer() {
	he := &httpEndpoint{
		paymentTopupUC: in.PaymentTopupUC,
	}

	in.Router.Endpoint(http.MethodPost, "/payments/topup", he.PaymentTopup)
}

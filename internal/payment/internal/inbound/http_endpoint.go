package inbound

import (
	"encoding/json"

	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shopspring/decimal"
)

var errInvalidBody = goerror.NewInvalidFormat("invalid request body")

type httpEndpoint struct {
	paymentTopupUC domain.PaymentTopup
}

func (e *httpEndpoint) PaymentTopup(c framework.Context) (any, error) {
	var req PaymentTopupRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	amo, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, errInvalidBody
	}

	resp, err := e.paymentTopupUC.Call(c.Context(), domain.PaymentTopupInput{
		ReferenceID: req.ReferenceID,
		Amount:      amo,
	})
	if err != nil {
		return nil, err
	}

	return PaymentTopupResponse{
		ReferenceID: resp.ReferenceID,
		Amount:      resp.Amount.StringFixed(2),
		Balance:     resp.Balance.StringFixed(2),
	}, nil
}

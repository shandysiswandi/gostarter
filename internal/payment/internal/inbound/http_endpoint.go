package inbound

import (
	"encoding/json"

	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shopspring/decimal"
)

var errInvalidBody = goerror.NewInvalidFormat("Request payload malformed")

type httpEndpoint struct {
	tel *telemetry.Telemetry

	paymentTopupUC domain.PaymentTopup
}

func (h *httpEndpoint) PaymentTopup(c framework.Context) (any, error) {
	ctx, span := h.tel.Tracer().Start(c.Context(), "payment.inbound.httpEndpoint.PaymentTopup")
	defer span.End()

	var req PaymentTopupRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	amo, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.paymentTopupUC.Call(ctx, domain.PaymentTopupInput{
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

package inbound

type (
	PaymentTopupRequest struct {
		ReferenceID string `json:"reference_id"`
		Amount      string `json:"amount"`
	}

	PaymentTopupResponse struct {
		ReferenceID string `json:"reference_id"`
		Amount      string `json:"amount"`
		Balance     string `json:"balance"`
	}
)

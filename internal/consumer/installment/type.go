package installment

type CreateInstallmentPayload struct {
	UserID         uint    `json:"user_id"`
	CreditOptionID uint    `json:"credit_limit_id"`
	OtrAmount      float64 `json:"otr_amount"`
	AssetName      string  `json:"asset_name"`
}

type PayInstallmentPayload struct {
	UserID     uint    `json:"user_id"`
	ContractID string  `json:"contract_id"`
	PaidAmount float64 `json:"paid_amount"`
}

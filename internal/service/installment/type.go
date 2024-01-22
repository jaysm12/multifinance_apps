package installment

import "errors"

// list Service error
var (
	ErrDataNotFound              = errors.New("data not found")
	ErrOtrAmountIsGreater        = errors.New("otr amount is greater than credit limit amount")
	ErrInvalidAmount             = errors.New("amount invalid")
	ErrInstallmentAlreadySettled = errors.New("installment already settled")
	ErrInvalidStatus             = errors.New("invalid installment status")
)

type CreateInstallmentRequest struct {
	CreditLimitID uint
	UserID        int
	OtrAmount     float64
	AssetName     string
}

type PayInstallmentRequest struct {
	InstallmentID uint
	PaidAmount    float64
}

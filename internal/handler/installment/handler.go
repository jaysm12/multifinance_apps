package installment

import "github.com/jaysm12/multifinance-apps/internal/service/installment"

type InstallmentHandler struct {
	service      installment.InstallmentServiceMethod
	timeoutInSec int
}

type Option func(*InstallmentHandler)

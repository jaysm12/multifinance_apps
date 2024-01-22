package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type InstallmentPaymentHistory struct {
	gorm.Model
	InstallmentID     uint
	ContractID        string
	InstallmentNumber int
	PaymentDate       time.Time
	PaidAmount        int
	RemainingAmount   int
}

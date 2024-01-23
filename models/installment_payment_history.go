package models

import (
	"time"

	"gorm.io/gorm"
)

type InstallmentPaymentHistory struct {
	gorm.Model
	InstallmentID     uint
	ContractID        string
	InstallmentNumber int
	PaymentDate       time.Time
	PaidAmount        float64
}

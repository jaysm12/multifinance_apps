package models

import "gorm.io/gorm"

type CreditLimit struct {
	gorm.Model
	UserID        uint
	DefaultAmount float64
	CurrentAmount float64
	Tenor         int
	Installments  []Installment
}

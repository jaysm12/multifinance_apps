package models

import "github.com/jinzhu/gorm"

type CreditLimit struct {
	gorm.Model
	UserID       uint
	Amount       int
	Tenor        int
	Installments []Installment
}

package models

import "gorm.io/gorm"

// User struct to user information
type UserKYC struct {
	gorm.Model
	UserID         uint
	NIK            string `gorm:"size:16;not null"`
	LegalName      string
	BirthDate      string
	BirthAddress   string
	SalaryAmount   string
	PhotoIDUrl     string
	PhotoSelfieUrl string
}

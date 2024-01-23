package util

import (
	"log"

	"github.com/jaysm12/multifinance-apps/models"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.UserKYC{},
		&models.CreditOption{},
		&models.Installment{},
		&models.InstallmentPaymentHistory{},
	)

	log.Println("Migration Done")
}

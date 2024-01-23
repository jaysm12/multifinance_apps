package installmentPaymentHistory

import (
	"errors"

	"github.com/jaysm12/multifinance-apps/models"
	"github.com/jaysm12/multifinance-apps/pkg/mysql"
	"gorm.io/gorm"
)

// InstallmentPaymentHistoryStoreMethod is set of methods for interacting with a installmentPaymentHistory storage system
type InstallmentPaymentHistoryStoreMethod interface {
	CreateInstallmentPaymentHistory(installmentPaymentHistoryinfo models.InstallmentPaymentHistory) error
	GetInstallmentPaymentHistoryInfoByID(installmentPaymentHistoryid int) (models.InstallmentPaymentHistory, error)
	GetLatestHistoryByInstallmentId(installmentID uint) (models.InstallmentPaymentHistory, error)
}

// InstallmentPaymentHistoryStore is list dependencies installmentPaymentHistory store
type InstallmentPaymentHistoryStore struct {
	db mysql.MysqlMethod
}

// NewInstallmentPaymentHistoryStore is func to generate InstallmentPaymentHistoryStoreMethod interface
func NewInstallmentPaymentHistoryStore(db mysql.MysqlMethod) InstallmentPaymentHistoryStoreMethod {
	return &InstallmentPaymentHistoryStore{
		db: db,
	}
}

func (u *InstallmentPaymentHistoryStore) getDB() (*gorm.DB, error) {
	db := u.db.GetDB()
	if db == nil {
		return nil, errors.New("database Client is not init")
	}

	return db, nil
}

// CreateInstallmentPaymentHistory is func to store / create installmentPaymentHistory info into database
func (u *InstallmentPaymentHistoryStore) CreateInstallmentPaymentHistory(installmentPaymentHistoryinfo models.InstallmentPaymentHistory) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	return db.Create(&installmentPaymentHistoryinfo).Error
}

// GetInstallmentPaymentHistoryByID is func to get installmentPaymentHistory info by id on database
func (u *InstallmentPaymentHistoryStore) GetInstallmentPaymentHistoryInfoByID(installmentPaymentHistoryid int) (models.InstallmentPaymentHistory, error) {
	var installmentPaymentHistory models.InstallmentPaymentHistory
	db, err := u.getDB()
	if err != nil {
		return models.InstallmentPaymentHistory{}, err
	}

	if err := db.First(&installmentPaymentHistory, installmentPaymentHistoryid).Error; err != nil {
		return models.InstallmentPaymentHistory{}, err
	}

	return installmentPaymentHistory, nil
}

func (u *InstallmentPaymentHistoryStore) GetLatestHistoryByInstallmentId(installmentID uint) (models.InstallmentPaymentHistory, error) {
	var installmentPaymentHistory models.InstallmentPaymentHistory
	db, err := u.getDB()
	if err != nil {
		return models.InstallmentPaymentHistory{}, err
	}

	if err := db.Where("installment_id = ?", installmentID).Last(&installmentPaymentHistory).Error; err != nil {
		return models.InstallmentPaymentHistory{}, err
	}

	return installmentPaymentHistory, nil
}

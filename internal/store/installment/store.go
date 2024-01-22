package installment

import (
	"errors"

	"github.com/jaysm12/multifinance-apps/models"
	"github.com/jaysm12/multifinance-apps/pkg/mysql"
	"github.com/jinzhu/gorm"
)

// InstallmentStoreMethod is set of methods for interacting with a installment storage system
type InstallmentStoreMethod interface {
	CreateInstallment(installmentinfo models.Installment) error
	GetInstallmentInfoByID(installmentid uint) (models.Installment, error)
	UpdateInstallment(installmentinfo models.Installment) error
}

// InstallmentStore is list dependencies installment store
type InstallmentStore struct {
	db mysql.MysqlMethod
}

// NewInstallmentStore is func to generate InstallmentStoreMethod interface
func NewInstallmentStore(db mysql.MysqlMethod) InstallmentStoreMethod {
	return &InstallmentStore{
		db: db,
	}
}

func (u *InstallmentStore) getDB() (*gorm.DB, error) {
	db := u.db.GetDB()
	if db == nil {
		return nil, errors.New("database Client is not init")
	}

	return db, nil
}

// CreateInstallment is func to store / create installment info into database
func (u *InstallmentStore) CreateInstallment(installmentinfo models.Installment) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	return db.Create(&installmentinfo).Error
}

// UpdateInstallment is func to edit / update installment info into database
func (u *InstallmentStore) UpdateInstallment(installmentinfo models.Installment) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	result := db.Model(models.Installment{}).Where("id = ?", installmentinfo.ID).Updates(installmentinfo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *InstallmentStore) GetInstallmentInfoByID(installmentid uint) (models.Installment, error) {
	var installment models.Installment
	db, err := u.getDB()
	if err != nil {
		return models.Installment{}, err
	}

	if err := db.First(&installment, installmentid).Error; err != nil {
		return models.Installment{}, err
	}

	return installment, nil
}

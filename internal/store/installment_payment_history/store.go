package installmentPaymentHistoryPaymentHistory

import (
	"errors"

	"github.com/jaysm12/multifinance-apps/models"
	"github.com/jaysm12/multifinance-apps/pkg/mysql"
	"github.com/jinzhu/gorm"
)

// InstallmentPaymentHistoryStoreMethod is set of methods for interacting with a installmentPaymentHistory storage system
type InstallmentPaymentHistoryStoreMethod interface {
	CreateInstallmentPaymentHistory(installmentPaymentHistoryinfo models.InstallmentPaymentHistory) error
	UpdateInstallmentPaymentHistory(installmentPaymentHistoryinfo models.InstallmentPaymentHistory) error
	DeleteInstallmentPaymentHistory(installmentPaymentHistoryid int) error
	GetInstallmentPaymentHistoryInfoByInstallmentPaymentHistoryname(installmentPaymentHistoryname string) (models.InstallmentPaymentHistory, error)
	GetInstallmentPaymentHistoryInfoByID(installmentPaymentHistoryid int) (models.InstallmentPaymentHistory, error)
	Count() (int, error)
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

// UpdateInstallmentPaymentHistory is func to edit / update installmentPaymentHistory info into database
func (u *InstallmentPaymentHistoryStore) UpdateInstallmentPaymentHistory(installmentPaymentHistoryinfo models.InstallmentPaymentHistory) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	result := db.Model(models.InstallmentPaymentHistory{}).Where("id = ?", installmentPaymentHistoryinfo.ID).Updates(installmentPaymentHistoryinfo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetInstallmentPaymentHistoryID is func to get installmentPaymentHistory id by installmentPaymentHistoryname and password
func (u *InstallmentPaymentHistoryStore) GetInstallmentPaymentHistoryInfoByInstallmentPaymentHistoryname(installmentPaymentHistoryname string) (models.InstallmentPaymentHistory, error) {
	var installmentPaymentHistory models.InstallmentPaymentHistory
	db, err := u.getDB()
	if err != nil {
		return models.InstallmentPaymentHistory{}, err
	}

	if err := db.Where("installmentPaymentHistoryname = ?", installmentPaymentHistoryname).First(&installmentPaymentHistory).Error; err != nil {
		return models.InstallmentPaymentHistory{}, err
	}

	return installmentPaymentHistory, nil
}

// DeleteInstallmentPaymentHistory is func to delete installmentPaymentHistory info on database
func (u *InstallmentPaymentHistoryStore) DeleteInstallmentPaymentHistory(installmentPaymentHistoryid int) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	installmentPaymentHistory := models.InstallmentPaymentHistory{
		Model: gorm.Model{
			ID: uint(installmentPaymentHistoryid),
		},
	}

	return db.Delete(&installmentPaymentHistory).Error
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

// Count is func to get total installmentPaymentHistory on database
func (u *InstallmentPaymentHistoryStore) Count() (int, error) {
	var installmentPaymentHistory models.InstallmentPaymentHistory
	db, err := u.getDB()
	if err != nil {
		return 0, err
	}

	var count int
	if err := db.Model(&installmentPaymentHistory).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

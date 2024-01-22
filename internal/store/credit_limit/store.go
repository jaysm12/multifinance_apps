package creditLimit

import (
	"errors"

	"github.com/jaysm12/multifinance-apps/models"
	"github.com/jaysm12/multifinance-apps/pkg/mysql"
	"github.com/jinzhu/gorm"
)

// CreditLimitStoreMethod is set of methods for interacting with a creditLimit storage system
type CreditLimitStoreMethod interface {
	CreateCreditLimitBulk(creditLimitinfo []models.CreditLimit) error
	UpdateCreditLimit(creditLimitinfo models.CreditLimit) error
	DeleteCreditLimit(creditLimitid uint) error
	GetCreditLimitInfoByID(creditLimitid uint) (models.CreditLimit, error)
	Count() (int, error)
}

// CreditLimitStore is list dependencies creditLimit store
type CreditLimitStore struct {
	db mysql.MysqlMethod
}

// NewCreditLimitStore is func to generate CreditLimitStoreMethod interface
func NewCreditLimitStore(db mysql.MysqlMethod) CreditLimitStoreMethod {
	return &CreditLimitStore{
		db: db,
	}
}

func (u *CreditLimitStore) getDB() (*gorm.DB, error) {
	db := u.db.GetDB()
	if db == nil {
		return nil, errors.New("database Client is not init")
	}

	return db, nil
}

// CreateCreditLimitBulk is func to store / create creditLimit info into database
func (u *CreditLimitStore) CreateCreditLimitBulk(creditLimitinfo []models.CreditLimit) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	return db.Create(&creditLimitinfo).Error
}

// UpdateCreditLimit is func to edit / update creditLimit info into database
func (u *CreditLimitStore) UpdateCreditLimit(creditLimitinfo models.CreditLimit) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	result := db.Model(models.CreditLimit{}).Where("id = ?", creditLimitinfo.ID).Updates(creditLimitinfo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteCreditLimit is func to delete creditLimit info on database
func (u *CreditLimitStore) DeleteCreditLimit(creditLimitid uint) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	creditLimit := models.CreditLimit{
		Model: gorm.Model{
			ID: uint(creditLimitid),
		},
	}

	return db.Delete(&creditLimit).Error
}

// GetCreditLimitByID is func to get creditLimit info by id on database
func (u *CreditLimitStore) GetCreditLimitInfoByID(creditLimitid uint) (models.CreditLimit, error) {
	var creditLimit models.CreditLimit
	db, err := u.getDB()
	if err != nil {
		return models.CreditLimit{}, err
	}

	if err := db.First(&creditLimit, creditLimitid).Error; err != nil {
		return models.CreditLimit{}, err
	}

	return creditLimit, nil
}

// Count is func to get total creditLimit on database
func (u *CreditLimitStore) Count() (int, error) {
	var creditLimit models.CreditLimit
	db, err := u.getDB()
	if err != nil {
		return 0, err
	}

	var count int
	if err := db.Model(&creditLimit).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

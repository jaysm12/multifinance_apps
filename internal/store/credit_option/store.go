package creditOption

import (
	"errors"

	"github.com/jaysm12/multifinance-apps/models"
	"github.com/jaysm12/multifinance-apps/pkg/mysql"
	"gorm.io/gorm"
)

// CreditOptionStoreMethod is set of methods for interacting with a creditOption storage system
type CreditOptionStoreMethod interface {
	CreateCreditOptionBulk(creditOptioninfo []models.CreditOption) error
	UpdateCreditOption(creditOptioninfo models.CreditOption) error
	DeleteCreditOption(creditOptionid uint) error
	GetCreditOptionInfoByID(creditOptionid uint) (models.CreditOption, error)
}

// CreditOptionStore is list dependencies creditOption store
type CreditOptionStore struct {
	db mysql.MysqlMethod
}

// NewCreditOptionStore is func to generate CreditOptionStoreMethod interface
func NewCreditOptionStore(db mysql.MysqlMethod) CreditOptionStoreMethod {
	return &CreditOptionStore{
		db: db,
	}
}

func (u *CreditOptionStore) getDB() (*gorm.DB, error) {
	db := u.db.GetDB()
	if db == nil {
		return nil, errors.New("database Client is not init")
	}

	return db, nil
}

// CreateCreditOptionBulk is func to store / create creditOption info into database
func (u *CreditOptionStore) CreateCreditOptionBulk(creditOptioninfo []models.CreditOption) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	return db.Create(&creditOptioninfo).Error
}

// UpdateCreditOption is func to edit / update creditOption info into database
func (u *CreditOptionStore) UpdateCreditOption(creditOptioninfo models.CreditOption) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	result := db.Model(models.CreditOption{}).Where("id = ?", creditOptioninfo.ID).Updates(creditOptioninfo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteCreditOption is func to delete creditOption info on database
func (u *CreditOptionStore) DeleteCreditOption(creditOptionid uint) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	creditOption := models.CreditOption{
		Model: gorm.Model{
			ID: uint(creditOptionid),
		},
	}

	return db.Delete(&creditOption).Error
}

// GetCreditOptionByID is func to get creditOption info by id on database
func (u *CreditOptionStore) GetCreditOptionInfoByID(creditOptionid uint) (models.CreditOption, error) {
	var creditOption models.CreditOption
	db, err := u.getDB()
	if err != nil {
		return models.CreditOption{}, err
	}

	if err := db.First(&creditOption, creditOptionid).Error; err != nil {
		return models.CreditOption{}, err
	}

	return creditOption, nil
}

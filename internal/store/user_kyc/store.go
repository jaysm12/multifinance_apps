package userkyc

import (
	"errors"

	"github.com/jaysm12/multifinance-apps/models"
	"github.com/jaysm12/multifinance-apps/pkg/mysql"
	"github.com/jinzhu/gorm"
)

type UserKYCStoreMethod interface {
	CreateUserKYC(userinfo models.UserKYC) error
	UpdateUserKYC(userinfo models.UserKYC) error
	DeleteUserKYC(userKycId int) error
	GetUserKYCInfoByUserID(userid int) (models.UserKYC, error)
	Count() (int, error)
}

// UserKYCStore is list dependencies user store
type UserKYCStore struct {
	db mysql.MysqlMethod
}

// NewUserKYCStore is func to generate UserKYCStoreMethod interface
func NewUserKYCStore(db mysql.MysqlMethod) UserKYCStoreMethod {
	return &UserKYCStore{
		db: db,
	}
}

func (u *UserKYCStore) getDB() (*gorm.DB, error) {
	db := u.db.GetDB()
	if db == nil {
		return nil, errors.New("database Client is not init")
	}

	return db, nil
}

// CreateUserKYC is func to store / create user info into database
func (u *UserKYCStore) CreateUserKYC(userinfo models.UserKYC) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	return db.Create(&userinfo).Error
}

// UpdateUserKYC is func to edit / update user info into database
func (u *UserKYCStore) UpdateUserKYC(userinfo models.UserKYC) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	result := db.Model(models.UserKYC{}).Where("id = ?", userinfo.ID).Updates(userinfo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetUserKYCID is func to get user id by username and password
func (u *UserKYCStore) GetUserKYCInfoByUserKYCname(username string) (models.UserKYC, error) {
	var user models.UserKYC
	db, err := u.getDB()
	if err != nil {
		return models.UserKYC{}, err
	}

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return models.UserKYC{}, err
	}

	return user, nil
}

// DeleteUserKYC is func to delete user info on database
func (u *UserKYCStore) DeleteUserKYC(userKycid int) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	user := models.UserKYC{
		Model: gorm.Model{
			ID: uint(userKycid),
		},
	}

	return db.Delete(&user).Error
}

// GetUserKYCByID is func to get user info by id on database
func (u *UserKYCStore) GetUserKYCInfoByUserID(userid int) (models.UserKYC, error) {
	var user models.UserKYC
	db, err := u.getDB()
	if err != nil {
		return models.UserKYC{}, err
	}

	if err := db.First(&user, userid).Error; err != nil {
		return models.UserKYC{}, err
	}

	return user, nil
}

// Count is func to get total user on database
func (u *UserKYCStore) Count() (int, error) {
	var user models.UserKYC
	db, err := u.getDB()
	if err != nil {
		return 0, err
	}

	var count int
	if err := db.Model(&user).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

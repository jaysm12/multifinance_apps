package user

import (
	"errors"

	"github.com/jaysm12/multifinance-apps/models"
	"github.com/jaysm12/multifinance-apps/pkg/mysql"
	"github.com/jinzhu/gorm"
)

// UserStoreMethod is set of methods for interacting with a user storage system
type UserStoreMethod interface {
	CreateUser(userinfo models.User) error
	UpdateUser(userinfo models.User) error
	DeleteUser(userid int) error
	GetUserInfoByUsername(username string) (models.User, error)
	GetUserInfoByID(userid int) (models.User, error)
	Count() (int, error)
}

// UserStore is list dependencies user store
type UserStore struct {
	db mysql.MysqlMethod
}

// NewUserStore is func to generate UserStoreMethod interface
func NewUserStore(db mysql.MysqlMethod) UserStoreMethod {
	return &UserStore{
		db: db,
	}
}

func (u *UserStore) getDB() (*gorm.DB, error) {
	db := u.db.GetDB()
	if db == nil {
		return nil, errors.New("database Client is not init")
	}

	return db, nil
}

// CreateUser is func to store / create user info into database
func (u *UserStore) CreateUser(userinfo models.User) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	return db.Create(&userinfo).Error
}

// UpdateUser is func to edit / update user info into database
func (u *UserStore) UpdateUser(userinfo models.User) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	var user models.User
	err = db.Where("username = ? AND id = ?", userinfo.Username, userinfo.ID).First(&user).Error
	if err != nil {
		return err
	}

	user.Password = userinfo.Password
	user.Fullname = userinfo.Fullname
	user.Email = userinfo.Email
	user.IsVerified = userinfo.IsVerified

	return db.Save(&user).Error
}

// GetUserID is func to get user id by username and password
func (u *UserStore) GetUserInfoByUsername(username string) (models.User, error) {
	var user models.User
	db, err := u.getDB()
	if err != nil {
		return models.User{}, err
	}

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

// DeleteUser is func to delete user info on database
func (u *UserStore) DeleteUser(userid int) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	user := models.User{
		Model: gorm.Model{
			ID: uint(userid),
		},
	}

	return db.Delete(&user).Error
}

// GetUserByID is func to get user info by id on database
func (u *UserStore) GetUserInfoByID(userid int) (models.User, error) {
	var user models.User
	db, err := u.getDB()
	if err != nil {
		return models.User{}, err
	}

	if err := db.First(&user, userid).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Count is func to get total user on database
func (u *UserStore) Count() (int, error) {
	var user models.User
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

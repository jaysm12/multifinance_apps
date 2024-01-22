package user

import (
	creditLimit "github.com/jaysm12/multifinance-apps/internal/store/credit_limit"
	"github.com/jaysm12/multifinance-apps/internal/store/user"
	userkyc "github.com/jaysm12/multifinance-apps/internal/store/user_kyc"
	"github.com/jaysm12/multifinance-apps/models"
	"github.com/jaysm12/multifinance-apps/pkg/hash"
)

// UserServiceMethod is list method for User Service
type UserServiceMethod interface {
	DeleteUser(DeleteUserServiceRequest) error
	UpdateUser(UpdateUserServiceRequest) (UserServiceInfo, error)
	GetUserByID(GetByIDServiceRequest) (UserServiceInfo, error)
	CreateUserKyc(request CreateUserKycRequest) error
}

// UserService is list dependencies for user service
type UserService struct {
	userStore        user.UserStoreMethod
	userKycStore     userkyc.UserKYCStoreMethod
	creditLimitStore creditLimit.CreditLimitStoreMethod
	hash             hash.HashMethod
}

// NewUserService is func to generate UserServiceMethod interface
func NewUserService(userStore user.UserStoreMethod, userKycStore userkyc.UserKYCStoreMethod, creditLimitStore creditLimit.CreditLimitStoreMethod, hash hash.HashMethod) UserServiceMethod {
	return &UserService{
		hash:             hash,
		userStore:        userStore,
		userKycStore:     userKycStore,
		creditLimitStore: creditLimitStore,
	}
}

// DeleteUser is service level func to validate and delete user info in database
func (u *UserService) DeleteUser(request DeleteUserServiceRequest) error {
	if request.UserId <= 0 {
		return ErrDataNotFound
	}

	userInfo, err := u.userStore.GetUserInfoByID(request.UserId)
	if err != nil || userInfo.ID <= 0 {
		return err
	}

	if !u.hash.CompareValue(userInfo.Password, request.Password) {
		return ErrPasswordIsIncorrect
	}

	return u.userStore.DeleteUser(int(userInfo.ID))
}

// UpdateUser is service level func to validate and update user info in database
func (u *UserService) UpdateUser(request UpdateUserServiceRequest) (UserServiceInfo, error) {
	if request.UserId <= 0 {
		return UserServiceInfo{}, ErrDataNotFound
	}

	userInfo, err := u.userStore.GetUserInfoByID(request.UserId)
	if err != nil || userInfo.ID <= 0 {
		return UserServiceInfo{}, err
	}

	if len(request.Password) > 0 {
		hashPassword, err := u.hash.HashValue(request.Password)
		if err != nil {
			return UserServiceInfo{}, err
		}

		userInfo.Password = string(hashPassword)
	}

	if len(request.Email) > 0 {
		userInfo.Email = request.Email
	}

	if len(request.Fullname) > 0 {
		userInfo.Fullname = request.Fullname
	}

	err = u.userStore.UpdateUser(userInfo)
	if err != nil {
		return UserServiceInfo{}, err
	}

	return UserServiceInfo{
		UserId:      int(userInfo.ID),
		Username:    userInfo.Username,
		Fullname:    userInfo.Fullname,
		Email:       userInfo.Email,
		CreatedDate: userInfo.CreatedAt.String(),
	}, nil
}

// GetUserByID is service level func to validate and get all user based id
func (u *UserService) GetUserByID(request GetByIDServiceRequest) (UserServiceInfo, error) {
	userInfo, err := u.userStore.GetUserInfoByID(int(request.UserId))
	if err != nil || userInfo.ID <= 0 {
		return UserServiceInfo{}, err
	}

	return UserServiceInfo{
		UserId:      int(userInfo.ID),
		Username:    userInfo.Username,
		Fullname:    userInfo.Fullname,
		Email:       userInfo.Email,
		IsVerified:  userInfo.IsVerified,
		CreatedDate: userInfo.CreatedAt.String(),
	}, nil
}

func (u *UserService) CreateUserKyc(request CreateUserKycRequest) error {
	userKycInfo := models.UserKYC{
		UserID:         uint(request.UserId),
		NIK:            request.NIK,
		LegalName:      request.LegalName,
		BirthDate:      request.BirthDate,
		BirthAddress:   request.BirthAddress,
		SalaryAmount:   request.SalaryAmount,
		PhotoIDUrl:     request.PhotoIDUrl,
		PhotoSelfieUrl: request.PhotoSelfieUrl,
	}

	// default credit limit
	creditInfo := []models.CreditLimit{
		{
			UserID: uint(request.UserId),
			Amount: 10000000,
			Tenor:  4,
		},
		{
			UserID: uint(request.UserId),
			Amount: 50000000,
			Tenor:  8,
		},
		{
			UserID: uint(request.UserId),
			Amount: 1000000,
			Tenor:  1,
		},
	}

	err := u.userKycStore.CreateUserKYC(userKycInfo)
	if err != nil {
		return err
	}

	err = u.creditLimitStore.CreateCreditLimitBulk(creditInfo)
	if err != nil {
		return err
	}

	return nil
}

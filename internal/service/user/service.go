package user

import (
	creditLimit "github.com/jaysm12/multifinance-apps/internal/store/credit_limit"
	"github.com/jaysm12/multifinance-apps/internal/store/user"
	userkyc "github.com/jaysm12/multifinance-apps/internal/store/user_kyc"
	"github.com/jaysm12/multifinance-apps/models"
)

// UserServiceMethod is list method for User Service
type UserServiceMethod interface {
	GetUserByID(GetByIDServiceRequest) (UserServiceInfo, error)
	CreateUserKyc(request CreateUserKycRequest) error
}

// UserService is list dependencies for user service
type UserService struct {
	userStore        user.UserStoreMethod
	userKycStore     userkyc.UserKYCStoreMethod
	creditLimitStore creditLimit.CreditLimitStoreMethod
}

// NewUserService is func to generate UserServiceMethod interface
func NewUserService(userStore user.UserStoreMethod, userKycStore userkyc.UserKYCStoreMethod, creditLimitStore creditLimit.CreditLimitStoreMethod) UserServiceMethod {
	return &UserService{
		userStore:        userStore,
		userKycStore:     userKycStore,
		creditLimitStore: creditLimitStore,
	}
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
			UserID:        uint(request.UserId),
			DefaultAmount: 100000000,
			CurrentAmount: 100000000,
			Tenor:         12,
		},
		{
			UserID:        uint(request.UserId),
			DefaultAmount: 500000000,
			CurrentAmount: 500000000,
			Tenor:         24,
		},
		{
			UserID:        uint(request.UserId),
			DefaultAmount: 800000000,
			CurrentAmount: 800000000,
			Tenor:         36,
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

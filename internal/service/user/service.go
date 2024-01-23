package user

import (
	"fmt"

	creditOption "github.com/jaysm12/multifinance-apps/internal/store/credit_option"
	"github.com/jaysm12/multifinance-apps/internal/store/user"
	userkyc "github.com/jaysm12/multifinance-apps/internal/store/user_kyc"
	"github.com/jaysm12/multifinance-apps/models"
)

// UserServiceMethod is list method for User Service
type UserServiceMethod interface {
	CreateUserKyc(request CreateUserKycRequest) error
}

// UserService is list dependencies for user service
type UserService struct {
	userStore         user.UserStoreMethod
	userKycStore      userkyc.UserKYCStoreMethod
	creditOptionStore creditOption.CreditOptionStoreMethod
}

// NewUserService is func to generate UserServiceMethod interface
func NewUserService(userStore user.UserStoreMethod, userKycStore userkyc.UserKYCStoreMethod, creditOptionStore creditOption.CreditOptionStoreMethod) UserServiceMethod {
	return &UserService{
		userStore:         userStore,
		userKycStore:      userKycStore,
		creditOptionStore: creditOptionStore,
	}
}

func (u *UserService) CreateUserKyc(request CreateUserKycRequest) error {
	userKycInfo := models.UserKYC{
		UserID:         request.UserId,
		NIK:            request.NIK,
		LegalName:      request.LegalName,
		BirthDate:      request.BirthDate,
		BirthAddress:   request.BirthAddress,
		SalaryAmount:   request.SalaryAmount,
		PhotoIDUrl:     request.PhotoIDUrl,
		PhotoSelfieUrl: request.PhotoSelfieUrl,
	}

	// default credit limit
	creditInfo := []models.CreditOption{
		{
			UserID:        request.UserId,
			DefaultAmount: 100000000,
			CurrentAmount: 100000000,
			Tenor:         12,
		},
		{
			UserID:        request.UserId,
			DefaultAmount: 500000000,
			CurrentAmount: 500000000,
			Tenor:         24,
		},
		{
			UserID:        request.UserId,
			DefaultAmount: 800000000,
			CurrentAmount: 800000000,
			Tenor:         36,
		},
	}

	err := u.userKycStore.CreateUserKYC(userKycInfo)
	if err != nil {
		return err
	}

	err = u.creditOptionStore.CreateCreditOptionBulk(creditInfo)
	if err != nil {
		fmt.Println("masuk")
		return err
	}

	return nil
}

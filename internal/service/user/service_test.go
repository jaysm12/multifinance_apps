package user

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	creditOption "github.com/jaysm12/multifinance-apps/internal/store/credit_option"
	mockCreditOption "github.com/jaysm12/multifinance-apps/internal/store/credit_option/mock"
	"github.com/jaysm12/multifinance-apps/internal/store/user"
	mockUser "github.com/jaysm12/multifinance-apps/internal/store/user/mock"
	userkyc "github.com/jaysm12/multifinance-apps/internal/store/user_kyc"
	mockUserKyc "github.com/jaysm12/multifinance-apps/internal/store/user_kyc/mock"
	"github.com/jaysm12/multifinance-apps/models"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		userStore         user.UserStoreMethod
		userKycStore      userkyc.UserKYCStoreMethod
		creditOptionStore creditOption.CreditOptionStoreMethod
	}
	tests := []struct {
		name string
		args args
		want UserServiceMethod
	}{
		{
			name: "success",
			args: args{
				userStore:         &user.UserStore{},
				userKycStore:      &userkyc.UserKYCStore{},
				creditOptionStore: &creditOption.CreditOptionStore{},
			},
			want: &UserService{
				userStore:         &user.UserStore{},
				userKycStore:      &userkyc.UserKYCStore{},
				creditOptionStore: &creditOption.CreditOptionStore{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.userStore, tt.args.userKycStore, tt.args.creditOptionStore); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_CreateUserKyc(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	uStore := mockUser.NewMockUserStoreMethod(mockCtrl)
	uKycStore := mockUserKyc.NewMockUserKYCStoreMethod(mockCtrl)
	cStore := mockCreditOption.NewMockCreditOptionStoreMethod(mockCtrl)
	defer mockCtrl.Finish()

	type args struct {
		request CreateUserKycRequest
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				request: CreateUserKycRequest{
					UserId:         1,
					NIK:            "1234567890123456",
					LegalName:      "John Doe",
					BirthDate:      "1990-05-15",
					BirthAddress:   "123 Main Street, Cityville",
					SalaryAmount:   "50000",
					PhotoIDUrl:     "http://example.com/photo_id.jpg",
					PhotoSelfieUrl: "http://example.com/photo_selfie.jpg",
				},
			},
			mockFunc: func() {
				uKycStore.EXPECT().CreateUserKYC(models.UserKYC{
					UserID:         1,
					NIK:            "1234567890123456",
					LegalName:      "John Doe",
					BirthDate:      "1990-05-15",
					BirthAddress:   "123 Main Street, Cityville",
					SalaryAmount:   "50000",
					PhotoIDUrl:     "http://example.com/photo_id.jpg",
					PhotoSelfieUrl: "http://example.com/photo_selfie.jpg",
				}).Return(nil)
				cStore.EXPECT().CreateCreditOptionBulk(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error kyc store flow",
			args: args{
				request: CreateUserKycRequest{
					UserId:         1,
					NIK:            "1234567890123456",
					LegalName:      "John Doe",
					BirthDate:      "1990-05-15",
					BirthAddress:   "123 Main Street, Cityville",
					SalaryAmount:   "50000",
					PhotoIDUrl:     "http://example.com/photo_id.jpg",
					PhotoSelfieUrl: "http://example.com/photo_selfie.jpg",
				},
			},
			mockFunc: func() {
				uKycStore.EXPECT().CreateUserKYC(
					models.UserKYC{
						UserID:         1,
						NIK:            "1234567890123456",
						LegalName:      "John Doe",
						BirthDate:      "1990-05-15",
						BirthAddress:   "123 Main Street, Cityville",
						SalaryAmount:   "50000",
						PhotoIDUrl:     "http://example.com/photo_id.jpg",
						PhotoSelfieUrl: "http://example.com/photo_selfie.jpg",
					}).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				userStore:         uStore,
				userKycStore:      uKycStore,
				creditOptionStore: cStore,
			}
			tt.mockFunc()
			if err := u.CreateUserKyc(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("UserService.CreateUserKyc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

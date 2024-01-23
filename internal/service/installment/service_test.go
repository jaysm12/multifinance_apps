package installment

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	creditOption "github.com/jaysm12/multifinance-apps/internal/store/credit_option"
	mockCreditOption "github.com/jaysm12/multifinance-apps/internal/store/credit_option/mock"
	"github.com/jaysm12/multifinance-apps/internal/store/installment"
	mockInstallment "github.com/jaysm12/multifinance-apps/internal/store/installment/mock"
	installmentPaymentHistory "github.com/jaysm12/multifinance-apps/internal/store/installment_payment_history"
	mockInstallmentPaymentHistory "github.com/jaysm12/multifinance-apps/internal/store/installment_payment_history/mock"
	"github.com/jaysm12/multifinance-apps/internal/store/user"
	mockUser "github.com/jaysm12/multifinance-apps/internal/store/user/mock"
	"github.com/jaysm12/multifinance-apps/models"
)

func TestNewInstallmentService(t *testing.T) {
	type args struct {
		storeInstallment               installment.InstallmentStoreMethod
		storeUser                      user.UserStoreMethod
		storeCreditOption              creditOption.CreditOptionStoreMethod
		storeInstallmentPaymentHistory installmentPaymentHistory.InstallmentPaymentHistoryStoreMethod
	}
	tests := []struct {
		name string
		args args
		want InstallmentServiceMethod
	}{
		{
			name: "success",
			args: args{
				storeInstallment:               &installment.InstallmentStore{},
				storeUser:                      &user.UserStore{},
				storeCreditOption:              &creditOption.CreditOptionStore{},
				storeInstallmentPaymentHistory: &installmentPaymentHistory.InstallmentPaymentHistoryStore{},
			},
			want: &InstallmentService{
				storeInstallment:               &installment.InstallmentStore{},
				storeUser:                      &user.UserStore{},
				storeCreditOption:              &creditOption.CreditOptionStore{},
				storeInstallmentPaymentHistory: &installmentPaymentHistory.InstallmentPaymentHistoryStore{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInstallmentService(tt.args.storeInstallment, tt.args.storeUser, tt.args.storeCreditOption, tt.args.storeInstallmentPaymentHistory); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstallmentService_CreateInstallment(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iStore := mockInstallment.NewMockInstallmentStoreMethod(mockCtrl)
	uStore := mockUser.NewMockUserStoreMethod(mockCtrl)
	cStore := mockCreditOption.NewMockCreditOptionStoreMethod(mockCtrl)
	iphStore := mockInstallmentPaymentHistory.NewMockInstallmentPaymentHistoryStoreMethod(mockCtrl)
	defer mockCtrl.Finish()

	type args struct {
		request CreateInstallmentRequest
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
				request: CreateInstallmentRequest{
					UserID:         1,
					CreditOptionID: 1,
					AssetName:      "Honda",
					OtrAmount:      20000000,
				},
			},
			mockFunc: func() {
				cl := models.CreditOption{
					CurrentAmount: 100000000,
					Tenor:         12,
				}
				cl.ID = uint(1)
				uStore.EXPECT().GetUserInfoByID(uint(1)).Return(models.User{}, nil)
				cStore.EXPECT().GetCreditOptionInfoByID(uint(1)).Return(cl, nil)
				iStore.EXPECT().CreateInstallment(gomock.Any()).Return(nil)
				cl.CurrentAmount = 80000000
				cStore.EXPECT().UpdateCreditOption(cl).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not enough limit",
			args: args{
				request: CreateInstallmentRequest{
					UserID:         1,
					CreditOptionID: 1,
					AssetName:      "Honda",
					OtrAmount:      20000000,
				},
			},
			mockFunc: func() {
				cl := models.CreditOption{
					CurrentAmount: 100000,
					Tenor:         12,
				}
				cl.ID = uint(1)
				uStore.EXPECT().GetUserInfoByID(uint(1)).Return(models.User{}, nil)
				cStore.EXPECT().GetCreditOptionInfoByID(uint(1)).Return(cl, nil)
			},
			wantErr: true,
		},
		{
			name: "error service",
			args: args{
				request: CreateInstallmentRequest{
					UserID:         1,
					CreditOptionID: 1,
					AssetName:      "Honda",
					OtrAmount:      20000000,
				},
			},
			mockFunc: func() {
				cl := models.CreditOption{
					CurrentAmount: 100000000,
					Tenor:         12,
				}
				cl.ID = uint(1)
				uStore.EXPECT().GetUserInfoByID(uint(1)).Return(models.User{}, nil)
				cStore.EXPECT().GetCreditOptionInfoByID(uint(1)).Return(cl, nil)
				iStore.EXPECT().CreateInstallment(gomock.Any()).Return(fmt.Errorf("some err"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InstallmentService{
				storeInstallment:               iStore,
				storeUser:                      uStore,
				storeCreditOption:              cStore,
				storeInstallmentPaymentHistory: iphStore,
			}
			tt.mockFunc()
			defer monkey.UnpatchAll()
			if err := i.CreateInstallment(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("InstallmentService.CreateInstallment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallmentService_PayInstallment(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iStore := mockInstallment.NewMockInstallmentStoreMethod(mockCtrl)
	uStore := mockUser.NewMockUserStoreMethod(mockCtrl)
	cStore := mockCreditOption.NewMockCreditOptionStoreMethod(mockCtrl)
	iphStore := mockInstallmentPaymentHistory.NewMockInstallmentPaymentHistoryStoreMethod(mockCtrl)
	defer mockCtrl.Finish()

	type args struct {
		request PayInstallmentRequest
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				request: PayInstallmentRequest{
					ContractID: "CN-1",
					PaidAmount: 1000,
					UserID:     1,
				},
			},
			mockFunc: func() {
				monkey.Patch(time.Now, func() time.Time { return time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC) })
				ins := models.Installment{
					Status:          models.InstallmentStatusInProgress,
					MonthlyAmount:   1000,
					ContractID:      "CN-1",
					RemainingAmount: 2000,
					CreditOptionID:  1,
					OtrAmount:       5000,
				}
				ins.ID = uint(1)
				iStore.EXPECT().GetInstallmentInfoByContractId("CN-1").Return(ins, nil)
				iphStore.EXPECT().GetLatestHistoryByInstallmentId(uint(1)).Return(models.InstallmentPaymentHistory{}, nil)

				ph := models.InstallmentPaymentHistory{
					InstallmentID:     uint(1),
					ContractID:        "CN-1",
					PaymentDate:       time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
					PaidAmount:        1000,
					InstallmentNumber: 1,
				}
				iphStore.EXPECT().CreateInstallmentPaymentHistory(ph).Return(nil)

				ins.RemainingAmount = 1000
				iStore.EXPECT().UpdateInstallment(ins).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success settled flow",
			args: args{
				request: PayInstallmentRequest{
					ContractID: "CN-1",
					PaidAmount: 1000,
					UserID:     1,
				},
			},
			mockFunc: func() {
				monkey.Patch(time.Now, func() time.Time { return time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC) })
				ins := models.Installment{
					Status:          models.InstallmentStatusInProgress,
					MonthlyAmount:   1000,
					ContractID:      "CN-1",
					RemainingAmount: 1000,
					CreditOptionID:  1,
					OtrAmount:       5000,
				}
				ins.ID = uint(1)
				iStore.EXPECT().GetInstallmentInfoByContractId("CN-1").Return(ins, nil)
				iphStore.EXPECT().GetLatestHistoryByInstallmentId(uint(1)).Return(models.InstallmentPaymentHistory{}, nil)

				ph := models.InstallmentPaymentHistory{
					InstallmentID:     uint(1),
					ContractID:        "CN-1",
					PaymentDate:       time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
					PaidAmount:        1000,
					InstallmentNumber: 1,
				}
				iphStore.EXPECT().CreateInstallmentPaymentHistory(ph).Return(nil)
				cl := models.CreditOption{
					CurrentAmount: 1000,
				}
				cStore.EXPECT().GetCreditOptionInfoByID(uint(1)).Return(cl, nil)
				cl.CurrentAmount = 6000
				cStore.EXPECT().UpdateCreditOption(cl).Return(nil)
				ins.RemainingAmount = 0
				ins.Status = models.InstallmentStatusSettled

				iStore.EXPECT().UpdateInstallment(ins).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error status invalid",
			args: args{
				request: PayInstallmentRequest{
					ContractID: "CN-1",
					PaidAmount: 1000,
					UserID:     1,
				},
			},
			mockFunc: func() {
				monkey.Patch(time.Now, func() time.Time { return time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC) })
				ins := models.Installment{
					Status:          models.InstallmentStatusFailed,
					MonthlyAmount:   1000,
					ContractID:      "CN-1",
					RemainingAmount: 2000,
					CreditOptionID:  1,
					OtrAmount:       5000,
				}
				ins.ID = uint(1)
				iStore.EXPECT().GetInstallmentInfoByContractId("CN-1").Return(ins, nil)
			},
			wantErr: true,
		},
		{
			name: "error status already settled",
			args: args{
				request: PayInstallmentRequest{
					ContractID: "CN-1",
					PaidAmount: 1000,
					UserID:     1,
				},
			},
			mockFunc: func() {
				monkey.Patch(time.Now, func() time.Time { return time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC) })
				ins := models.Installment{
					Status:          models.InstallmentStatusSettled,
					MonthlyAmount:   1000,
					ContractID:      "CN-1",
					RemainingAmount: 0,
					CreditOptionID:  1,
					OtrAmount:       5000,
				}
				ins.ID = uint(1)
				iStore.EXPECT().GetInstallmentInfoByContractId("CN-1").Return(ins, nil)
			},
			wantErr: true,
		},
		{
			name: "error service error",
			args: args{
				request: PayInstallmentRequest{
					ContractID: "CN-1",
					PaidAmount: 1000,
					UserID:     1,
				},
			},
			mockFunc: func() {
				iStore.EXPECT().GetInstallmentInfoByContractId("CN-1").Return(models.Installment{}, fmt.Errorf("some err"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InstallmentService{
				storeInstallment:               iStore,
				storeUser:                      uStore,
				storeCreditOption:              cStore,
				storeInstallmentPaymentHistory: iphStore,
			}
			tt.mockFunc()
			defer monkey.UnpatchAll()
			if err := i.PayInstallment(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("InstallmentService.PayInstallment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

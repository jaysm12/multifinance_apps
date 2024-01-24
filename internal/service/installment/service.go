package installment

import (
	"fmt"
	"math/rand"
	"time"

	creditOption "github.com/jaysm12/multifinance-apps/internal/store/credit_option"
	"github.com/jaysm12/multifinance-apps/internal/store/installment"
	installmentPaymentHistory "github.com/jaysm12/multifinance-apps/internal/store/installment_payment_history"
	"github.com/jaysm12/multifinance-apps/internal/store/user"
	"github.com/jaysm12/multifinance-apps/models"
	"gorm.io/gorm"
)

type InstallmentServiceMethod interface {
	CreateInstallment(request CreateInstallmentRequest) error
	PayInstallment(request PayInstallmentRequest) error
}

type InstallmentService struct {
	storeInstallment               installment.InstallmentStoreMethod
	storeUser                      user.UserStoreMethod
	storeCreditOption              creditOption.CreditOptionStoreMethod
	storeInstallmentPaymentHistory installmentPaymentHistory.InstallmentPaymentHistoryStoreMethod
}

func NewInstallmentService(
	storeInstallment installment.InstallmentStoreMethod,
	storeUser user.UserStoreMethod,
	storeCreditOption creditOption.CreditOptionStoreMethod,
	storeInstallmentPaymentHistory installmentPaymentHistory.InstallmentPaymentHistoryStoreMethod,
) InstallmentServiceMethod {
	return &InstallmentService{
		storeInstallment:               storeInstallment,
		storeUser:                      storeUser,
		storeCreditOption:              storeCreditOption,
		storeInstallmentPaymentHistory: storeInstallmentPaymentHistory,
	}
}

const (
	defaultInterest = 5.0
	defaultAdminFee = 100000.0
)

func (i *InstallmentService) CreateInstallment(request CreateInstallmentRequest) error {
	_, err := i.storeUser.GetUserInfoByID(request.UserID)
	if err != nil {
		return ErrDataNotFound
	}

	cl, err := i.storeCreditOption.GetCreditOptionInfoByID(request.CreditOptionID)
	if err != nil {
		return ErrDataNotFound
	}

	// error if limit is not enough
	if request.OtrAmount > cl.CurrentAmount {
		return ErrOtrAmountIsGreater
	}
	adminFee := defaultAdminFee

	totalInstallment, monthlyAmount, totalInterest, interestPerMonth := calculateInstallmentDetails(request.OtrAmount+adminFee, defaultInterest, cl.Tenor)

	installmentInfo := models.Installment{
		UserID:                 uint(request.UserID),
		CreditOptionID:         uint(request.CreditOptionID),
		ContractID:             generateContractID(),
		AssetName:              request.AssetName,
		OtrAmount:              request.OtrAmount,
		InterestRate:           defaultInterest,
		TotalInstallmentAmount: totalInstallment,
		MonthlyAmount:          monthlyAmount,
		TotalInterest:          totalInterest,
		InterestPerMonth:       interestPerMonth,
		Status:                 models.InstallmentStatusInProgress,
		AdminFee:               defaultAdminFee,
		Tenor:                  cl.Tenor,
		RemainingAmount:        totalInstallment,
	}

	err = i.storeInstallment.CreateInstallment(installmentInfo)
	if err != nil {
		return err
	}

	cl.CurrentAmount = cl.CurrentAmount - request.OtrAmount
	err = i.storeCreditOption.UpdateCreditOption(cl)
	if err != nil {
		return err
	}

	return nil
}

func (i *InstallmentService) PayInstallment(request PayInstallmentRequest) error {
	installmentInfo, err := i.storeInstallment.GetInstallmentInfoByContractId(request.ContractID)
	if err != nil {
		return ErrDataNotFound
	}
	if installmentInfo.Status != models.InstallmentStatusInProgress && installmentInfo.Status != models.InstallmentStatusOverdue {
		if installmentInfo.Status == models.InstallmentStatusSettled {
			return ErrInstallmentAlreadySettled
		}
		return ErrInvalidStatus
	}

	if request.PaidAmount != installmentInfo.MonthlyAmount {
		return ErrInvalidAmount
	}

	latestHistory, err := i.storeInstallmentPaymentHistory.GetLatestHistoryByInstallmentId(installmentInfo.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	installmentNumber := 1
	if latestHistory.ID > 0 {
		installmentNumber = latestHistory.InstallmentNumber + 1
	}
	paymentHistory := models.InstallmentPaymentHistory{
		InstallmentID:     installmentInfo.ID,
		ContractID:        installmentInfo.ContractID,
		PaymentDate:       time.Now(),
		PaidAmount:        request.PaidAmount,
		InstallmentNumber: installmentNumber,
	}

	err = i.storeInstallmentPaymentHistory.CreateInstallmentPaymentHistory(paymentHistory)
	if err != nil {
		return err
	}

	remainingAmount := installmentInfo.RemainingAmount - request.PaidAmount

	// if installment is settled
	if remainingAmount == 0 {
		installmentInfo.Status = models.InstallmentStatusSettled

		cl, err := i.storeCreditOption.GetCreditOptionInfoByID(installmentInfo.CreditOptionID)
		if err != nil {
			return err
		}

		cl.CurrentAmount = cl.CurrentAmount + installmentInfo.OtrAmount
		err = i.storeCreditOption.UpdateCreditOption(cl)
		if err != nil {
			return err
		}
	}

	installmentInfo.RemainingAmount = remainingAmount
	err = i.storeInstallment.UpdateInstallment(installmentInfo)
	if err != nil {
		return err
	}

	return nil
}

func generateContractID() string {
	// Get the current date and time
	currentTime := time.Now()

	// Generate a random string
	rand.Seed(currentTime.UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 7
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	contractID := fmt.Sprintf("CN-%v%v", string(result), currentTime.Format("20060102150405"))

	return contractID
}

func calculateInstallmentDetails(amount, interest float64, tenorMonths int) (float64, float64, float64, float64) {
	// Convert fixed interest rate to decimal
	r := interest / 100

	totalInterest := amount * r * float64(tenorMonths)
	totalInstallment := amount + totalInterest
	monthlyAmount := totalInstallment / float64(tenorMonths)
	monthlyInterest := totalInterest / float64(tenorMonths)

	return totalInstallment, monthlyAmount, totalInterest, monthlyInterest
}

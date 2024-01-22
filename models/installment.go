package models

import (
	"github.com/jinzhu/gorm"
)

type InstallmentStatus int

const (
	InstallmentStatusUnknown    InstallmentStatus = -1
	InstallmentStatusInProgress InstallmentStatus = 1
	InstallmentStatusSettled    InstallmentStatus = 2
	InstallmentStatusCancelled  InstallmentStatus = 3
	InstallmentStatusOverdue    InstallmentStatus = 3
	InstallmentStatusFailed     InstallmentStatus = 4
)

type Installment struct {
	gorm.Model
	CreditLimitID          uint
	ContractID             string
	AssetName              string
	OtrAmount              float64
	AdminFee               int
	TotalInterest          float64
	MonthlyAmount          float64
	TotalInstallmentAmount float64
	InterestRate           float64
	InterestPerMonth       float64
	Status                 InstallmentStatus
	OverdueAmount          int
	OverdueDays            int
	Tenor                  int
	RemainingAmount        float64
}

var InstallmentStatusSelectorString = map[int]string{
	int(InstallmentStatusUnknown):    "UNKNOWN",
	int(InstallmentStatusInProgress): "INPROGRESS",
	int(InstallmentStatusSettled):    "SETTLED",
	int(InstallmentStatusFailed):     "FAILED",
}

func (v InstallmentStatus) String() string {
	if str, ok := InstallmentStatusSelectorString[int(v)]; ok {
		return str
	}
	return InstallmentStatusSelectorString[-1]
}

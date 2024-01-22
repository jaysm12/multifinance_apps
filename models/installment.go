package models

import "github.com/jinzhu/gorm"

type InstallmentStatus int

const (
	InstallmentStatusUnknown    InstallmentStatus = -1
	InstallmentStatusInProgress InstallmentStatus = 1
	InstallmentStatusSettled    InstallmentStatus = 2
	InstallmentStatusFailed     InstallmentStatus = 3
)

type Installment struct {
	gorm.Model
	ContractID        string
	AssetName         string
	OtrAmount         int
	AdminFee          int
	InterestAmount    int
	InstallmentAmount int
	Status            InstallmentStatus
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

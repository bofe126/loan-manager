package models

import (
	"time"
)

// Loan 贷款模型
type Loan struct {
	ID             uint          `gorm:"primaryKey" json:"id"`
	BorrowerName   string        `gorm:"size:100" json:"borrower_name"`
	BankName       string        `gorm:"size:100" json:"bank_name"`
	Amount         float64       `json:"amount"`
	InterestRate   float64       `json:"interest_rate"`
	StartDate      time.Time     `json:"start_date"`
	EndDate        time.Time     `json:"end_date"`
	PaymentDate    int           `gorm:"default:1" json:"payment_date"`
	Status         LoanStatus    `gorm:"size:20;default:active" json:"status"`
	PaymentMethod  PaymentMethod `gorm:"size:50;default:等额本息" json:"payment_method"`
	MonthlyPayment float64       `json:"monthly_payment"`
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"created_at"`
	Payments       []Payment     `gorm:"foreignKey:LoanID" json:"payments,omitempty"`
}

// TableName 指定表名
func (Loan) TableName() string {
	return "loan"
}

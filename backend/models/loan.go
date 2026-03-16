package models

import (
	"time"
)

// Loan 贷款模型
type Loan struct {
	ID                 uint          `gorm:"primaryKey" json:"id"`
	BorrowerName       string        `gorm:"size:100" json:"borrower_name"`
	BankName           string        `gorm:"size:100" json:"bank_name"`
	TotalAmount        float64       `gorm:"column:total_amount" json:"total_amount"`   // 贷款总额（原始金额）
	RemainingAmount    float64       `gorm:"column:remaining_amount" json:"remaining_amount"` // 当前余额
	InterestRate       float64       `json:"interest_rate"`
	StartDate          time.Time     `json:"start_date"`                                // 贷款开始日期
	EndDate            time.Time     `json:"end_date"`                                  // 贷款结束日期
	FirstPaymentDate   *time.Time    `json:"first_payment_date,omitempty"`              // 第一期还款日期（可选）
	PaymentDate        int           `gorm:"default:1" json:"payment_date"`             // 每月还款日（1-31）
	Status             LoanStatus    `gorm:"size:20;default:active" json:"status"`
	PaymentMethod      PaymentMethod `gorm:"size:50;default:等额本息" json:"payment_method"`
	MonthlyPayment     float64       `json:"monthly_payment"`
	CreatedAt          time.Time     `gorm:"autoCreateTime" json:"created_at"`
	Payments           []Payment     `gorm:"foreignKey:LoanID" json:"payments,omitempty"`
}

// TableName 指定表名
func (Loan) TableName() string {
	return "loan"
}

package models

import (
	"time"
)

// Payment 还款记录模型
type Payment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LoanID      uint      `gorm:"not null;index" json:"loan_id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	PaymentDate time.Time `gorm:"autoCreateTime" json:"payment_date"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payment"
}

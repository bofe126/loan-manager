package services

import (
	"loan-manager-wails/backend/models"
	"time"

	"gorm.io/gorm"
)

// PaymentService 还款服务
type PaymentService struct {
	db                *gorm.DB
	calculatorService *CalculatorService
}

// NewPaymentService 创建还款服务
func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{
		db:                db,
		calculatorService: NewCalculatorService(),
	}
}

// MakePayment 记录还款
func (s *PaymentService) MakePayment(loanID uint, amount float64) error {
	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 查询贷款
		var loan models.Loan
		if err := tx.First(&loan, loanID).Error; err != nil {
			return err
		}

		// 创建还款记录
		payment := models.Payment{
			LoanID:      loanID,
			Amount:      amount,
			PaymentDate: time.Now(),
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		// 更新贷款金额
		loan.Amount -= amount
		if loan.Amount <= 0 {
			loan.Amount = 0
			loan.Status = models.StatusCompleted
		}

		// 重新计算月供
		if loan.Amount > 0 {
			loan.MonthlyPayment = s.calculatorService.CalculateMonthlyPayment(
				loan.Amount,
				loan.InterestRate,
				loan.StartDate,
				loan.EndDate,
				loan.PaymentMethod,
			)
		} else {
			loan.MonthlyPayment = 0
		}

		// 保存贷款
		return tx.Save(&loan).Error
	})
}

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
func (s *PaymentService) MakePayment(loanID uint, amount float64, paymentDate time.Time) error {
	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 查询贷款（预加载还款记录）
		var loan models.Loan
		if err := tx.Preload("Payments").First(&loan, loanID).Error; err != nil {
			return err
		}

		// 创建还款记录
		payment := models.Payment{
			LoanID:      loanID,
			Amount:      amount,
			PaymentDate: paymentDate,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		// 重新加载还款记录（包含新创建的）
		loan.Payments = append(loan.Payments, payment)

		// 计算新的剩余金额（考虑正常还款 + 额外还款）
		newRemaining := s.calculatorService.CalculateRemainingAmount(loan, paymentDate)

		// 更新贷款信息
		loan.RemainingAmount = newRemaining

		if newRemaining <= 0 {
			loan.Status = models.StatusCompleted
			loan.MonthlyPayment = 0
		} else {
			// 一次性还款后，需要重新计算月供
			// 使用新的剩余金额和剩余期限重新计算
			loan.MonthlyPayment = s.calculatorService.CalculateMonthlyPayment(
				newRemaining,
				loan.InterestRate,
				paymentDate,
				loan.EndDate,
				loan.PaymentMethod,
			)
		}

		// 保存贷款
		return tx.Save(&loan).Error
	})
}

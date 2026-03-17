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
			// 还款后，需要重新计算月供
			// 计算下一个还款日
			var nextPaymentDate time.Time
			if loan.FirstPaymentDate != nil {
				nextPaymentDate = *loan.FirstPaymentDate
			} else {
				nextPaymentDate = loan.StartDate.AddDate(0, 1, 0)
				nextPaymentDate = time.Date(
					nextPaymentDate.Year(),
					nextPaymentDate.Month(),
					loan.PaymentDate,
					0, 0, 0, 0,
					nextPaymentDate.Location(),
				)
			}

			// 找到还款日之后的下一个还款日
			for nextPaymentDate.Before(paymentDate) || nextPaymentDate.Equal(paymentDate) {
				nextPaymentDate = nextPaymentDate.AddDate(0, 1, 0)
			}

			// 计算最后还款日（基于还款日，而不是结束日期）
			lastPaymentDate := time.Date(
				loan.EndDate.Year(),
				loan.EndDate.Month(),
				loan.PaymentDate,
				0, 0, 0, 0,
				loan.EndDate.Location(),
			)
			// 如果最后还款日在结束日期之后，往前推一个月
			if lastPaymentDate.After(loan.EndDate) {
				lastPaymentDate = lastPaymentDate.AddDate(0, -1, 0)
			}

			loan.MonthlyPayment = s.calculatorService.CalculateMonthlyPayment(
				newRemaining,
				loan.InterestRate,
				nextPaymentDate,
				lastPaymentDate,
				loan.PaymentMethod,
			)
		}

		// 保存贷款
		return tx.Save(&loan).Error
	})
}

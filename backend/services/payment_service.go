package services

import (
	"fmt"
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

		fmt.Printf("\n=== MakePayment Debug ===\n")
		fmt.Printf("贷款ID: %d\n", loanID)
		fmt.Printf("还款金额: %.2f\n", amount)
		fmt.Printf("还款日期: %s\n", paymentDate.Format("2006-01-02"))
		fmt.Printf("还款前余额: %.2f\n", loan.RemainingAmount)
		fmt.Printf("还款前月供: %.2f\n", loan.MonthlyPayment)
		fmt.Printf("已有还款记录数: %d\n", len(loan.Payments))

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
		fmt.Printf("添加还款记录后，总记录数: %d\n", len(loan.Payments))

		// 计算新的剩余金额（考虑正常还款 + 额外还款）
		newRemaining := s.calculatorService.CalculateRemainingAmount(loan, paymentDate)
		fmt.Printf("计算出的新余额: %.2f\n", newRemaining)

		// 更新贷款信息
		loan.RemainingAmount = newRemaining

		if newRemaining <= 0 {
			loan.Status = models.StatusCompleted
			loan.MonthlyPayment = 0
			fmt.Printf("余额为0，设置状态为已完成\n")
		} else {
			// 还款后，需要重新计算月供
			fmt.Printf("=== 重新计算月供 ===\n")

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
			fmt.Printf("下一个还款日: %s\n", nextPaymentDate.Format("2006-01-02"))

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
			fmt.Printf("最后还款日: %s\n", lastPaymentDate.Format("2006-01-02"))

			loan.MonthlyPayment = s.calculatorService.CalculateMonthlyPayment(
				newRemaining,
				loan.InterestRate,
				nextPaymentDate,
				lastPaymentDate,
				loan.PaymentMethod,
			)
			fmt.Printf("计算出的新月供: %.2f\n", loan.MonthlyPayment)
		}

		fmt.Printf("保存到数据库: 余额=%.2f, 月供=%.2f\n", loan.RemainingAmount, loan.MonthlyPayment)
		fmt.Printf("=== MakePayment 完成 ===\n\n")

		// 保存贷款
		return tx.Save(&loan).Error
	})
}

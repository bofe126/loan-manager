package services

import (
	"loan-manager-wails/backend/models"
	"time"

	"gorm.io/gorm"
)

// LoanService 贷款服务
type LoanService struct {
	db                *gorm.DB
	calculatorService *CalculatorService
}

// NewLoanService 创建贷款服务
func NewLoanService(db *gorm.DB) *LoanService {
	return &LoanService{
		db:                db,
		calculatorService: NewCalculatorService(),
	}
}

// GetAllLoans 获取所有贷款
func (s *LoanService) GetAllLoans() ([]models.Loan, error) {
	var loans []models.Loan
	err := s.db.Preload("Payments").Order("created_at DESC").Find(&loans).Error
	if err != nil {
		return nil, err
	}

	// 实时计算每笔贷款的当前余额和月供
	now := time.Now()
	for i := range loans {
		loans[i].RemainingAmount = s.calculatorService.CalculateRemainingAmount(loans[i], now)

		// 如果余额为0，更新状态为已完成
		if loans[i].RemainingAmount <= 0 && loans[i].Status != models.StatusCompleted {
			loans[i].Status = models.StatusCompleted
			loans[i].MonthlyPayment = 0
			s.db.Model(&loans[i]).Updates(map[string]interface{}{
				"status":          models.StatusCompleted,
				"monthly_payment": 0,
			})
		} else if loans[i].RemainingAmount > 0 {
			// 找到最后一次实际还款日期
			var lastActualPaymentDate time.Time
			if len(loans[i].Payments) > 0 {
				for _, payment := range loans[i].Payments {
					if payment.PaymentDate.After(lastActualPaymentDate) {
						lastActualPaymentDate = payment.PaymentDate
					}
				}
			}

			// 如果有实际还款记录，基于最后还款日计算；否则基于开始日期
			var baseDate time.Time
			var baseBalance float64
			if !lastActualPaymentDate.IsZero() {
				baseDate = lastActualPaymentDate
				baseBalance = s.calculatorService.CalculateRemainingAmount(loans[i], lastActualPaymentDate)
			} else {
				baseDate = loans[i].StartDate
				baseBalance = loans[i].TotalAmount
			}

			// 计算下一个还款日
			var nextPaymentDate time.Time
			if loans[i].FirstPaymentDate != nil {
				nextPaymentDate = *loans[i].FirstPaymentDate
			} else {
				nextPaymentDate = loans[i].StartDate.AddDate(0, 1, 0)
				nextPaymentDate = time.Date(
					nextPaymentDate.Year(),
					nextPaymentDate.Month(),
					loans[i].PaymentDate,
					0, 0, 0, 0,
					nextPaymentDate.Location(),
				)
			}

			// 找到baseDate之后的下一个还款日
			for nextPaymentDate.Before(baseDate) || nextPaymentDate.Equal(baseDate) {
				nextPaymentDate = nextPaymentDate.AddDate(0, 1, 0)
			}

			// 计算最后还款日
			lastPaymentDate := time.Date(
				loans[i].EndDate.Year(),
				loans[i].EndDate.Month(),
				loans[i].PaymentDate,
				0, 0, 0, 0,
				loans[i].EndDate.Location(),
			)
			if lastPaymentDate.After(loans[i].EndDate) {
				lastPaymentDate = lastPaymentDate.AddDate(0, -1, 0)
			}

			// 计算剩余期数
			remainingMonths := s.calculatorService.CalculateLoanMonths(nextPaymentDate, lastPaymentDate) + 1
			monthlyRate := s.calculatorService.CalculateMonthlyRate(loans[i].InterestRate)
			loans[i].MonthlyPayment = s.calculatorService.CalculateMonthlyPaymentForAmount(
				baseBalance,
				monthlyRate,
				remainingMonths,
			)
		}
	}

	return loans, nil
}

// GetLoanByID 根据ID获取贷款
func (s *LoanService) GetLoanByID(id uint) (*models.Loan, error) {
	var loan models.Loan
	err := s.db.Preload("Payments").First(&loan, id).Error
	if err != nil {
		return nil, err
	}

	// 实时计算当前余额和月供
	now := time.Now()
	loan.RemainingAmount = s.calculatorService.CalculateRemainingAmount(loan, now)

	// 如果余额为0，更新状态为已完成
	if loan.RemainingAmount <= 0 && loan.Status != models.StatusCompleted {
		loan.Status = models.StatusCompleted
		loan.MonthlyPayment = 0
		s.db.Model(&loan).Updates(map[string]interface{}{
			"status":          models.StatusCompleted,
			"monthly_payment": 0,
		})
	} else if loan.RemainingAmount > 0 {
		// 找到最后一次实际还款日期
		var lastActualPaymentDate time.Time
		if len(loan.Payments) > 0 {
			for _, payment := range loan.Payments {
				if payment.PaymentDate.After(lastActualPaymentDate) {
					lastActualPaymentDate = payment.PaymentDate
				}
			}
		}

		// 如果有实际还款记录，基于最后还款日计算；否则基于开始日期
		var baseDate time.Time
		var baseBalance float64
		if !lastActualPaymentDate.IsZero() {
			baseDate = lastActualPaymentDate
			baseBalance = s.calculatorService.CalculateRemainingAmount(loan, lastActualPaymentDate)
		} else {
			baseDate = loan.StartDate
			baseBalance = loan.TotalAmount
		}

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

		// 找到baseDate之后的下一个还款日
		for nextPaymentDate.Before(baseDate) || nextPaymentDate.Equal(baseDate) {
			nextPaymentDate = nextPaymentDate.AddDate(0, 1, 0)
		}

		// 计算最后还款日
		lastPaymentDate := time.Date(
			loan.EndDate.Year(),
			loan.EndDate.Month(),
			loan.PaymentDate,
			0, 0, 0, 0,
			loan.EndDate.Location(),
		)
		if lastPaymentDate.After(loan.EndDate) {
			lastPaymentDate = lastPaymentDate.AddDate(0, -1, 0)
		}

		// 计算剩余期数
		remainingMonths := s.calculatorService.CalculateLoanMonths(nextPaymentDate, lastPaymentDate) + 1
		monthlyRate := s.calculatorService.CalculateMonthlyRate(loan.InterestRate)
		loan.MonthlyPayment = s.calculatorService.CalculateMonthlyPaymentForAmount(
			baseBalance,
			monthlyRate,
			remainingMonths,
		)
	}

	return &loan, nil
}

// CreateLoan 创建贷款
func (s *LoanService) CreateLoan(loan *models.Loan) error {
	// 初始化当前余额
	if loan.RemainingAmount == 0 {
		loan.RemainingAmount = loan.TotalAmount
	}

	// 计算月供
	loan.MonthlyPayment = s.calculatorService.CalculateMonthlyPayment(
		loan.TotalAmount,
		loan.InterestRate,
		loan.StartDate,
		loan.EndDate,
		loan.PaymentMethod,
	)

	// 设置默认状态
	if loan.Status == "" {
		loan.Status = models.StatusActive
	}

	return s.db.Create(loan).Error
}

// UpdateLoan 更新贷款
func (s *LoanService) UpdateLoan(loan *models.Loan) error {
	// 重新计算月供
	loan.MonthlyPayment = s.calculatorService.CalculateMonthlyPayment(
		loan.TotalAmount,
		loan.InterestRate,
		loan.StartDate,
		loan.EndDate,
		loan.PaymentMethod,
	)

	return s.db.Save(loan).Error
}

// DeleteLoan 删除贷款
func (s *LoanService) DeleteLoan(id uint) error {
	// 先删除关联的还款记录
	if err := s.db.Where("loan_id = ?", id).Delete(&models.Payment{}).Error; err != nil {
		return err
	}

	// 删除贷款
	return s.db.Delete(&models.Loan{}, id).Error
}

package services

import (
	"fmt"
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

	// 计算每笔贷款的当前余额和月供
	now := time.Now()
	for i := range loans {
		// 找到最后一次实际还款日期
		var lastActualPaymentDate time.Time
		if len(loans[i].Payments) > 0 {
			for _, payment := range loans[i].Payments {
				if payment.PaymentDate.After(lastActualPaymentDate) {
					lastActualPaymentDate = payment.PaymentDate
				}
			}
		}

		// 余额：模拟到当前日期
		loans[i].RemainingAmount = s.calculatorService.CalculateRemainingAmount(loans[i], now)

		// 如果余额为0，更新状态为已完成
		if loans[i].RemainingAmount <= 0 && loans[i].Status != models.StatusCompleted {
			loans[i].Status = models.StatusCompleted
			loans[i].MonthlyPayment = 0
			s.db.Model(&loans[i]).Updates(map[string]interface{}{
				"status":          models.StatusCompleted,
				"monthly_payment": 0,
			})
		} else {
			// 计算最后还款后的余额
			var baseBalance float64
			var baseDate time.Time
			if !lastActualPaymentDate.IsZero() {
				baseBalance = s.calculatorService.CalculateRemainingAmount(loans[i], lastActualPaymentDate.AddDate(0, 0, 1))
				baseDate = lastActualPaymentDate
			} else {
				baseBalance = loans[i].TotalAmount
				baseDate = loans[i].StartDate
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

			for nextPaymentDate.Before(baseDate) || nextPaymentDate.Equal(baseDate) {
				nextPaymentDate = nextPaymentDate.AddDate(0, 1, 0)
			}

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

			// 计算剩余期数（+1包含首尾）
			remainingMonths := s.calculatorService.CalculateLoanMonths(nextPaymentDate, lastPaymentDate) + 1
			monthlyRate := s.calculatorService.CalculateMonthlyRate(loans[i].InterestRate)

			// 根据还款方式计算月供
			switch loans[i].PaymentMethod {
			case models.EqualInstallment:
				// 等额本息
				if monthlyRate == 0 {
					loans[i].MonthlyPayment = baseBalance / float64(remainingMonths)
				} else {
					pow := 1.0
					for j := 0; j < remainingMonths; j++ {
						pow *= (1 + monthlyRate)
					}
					loans[i].MonthlyPayment = baseBalance * monthlyRate * pow / (pow - 1)
				}

			case models.EqualPrincipal:
				// 等额本金
				principalPerMonth := baseBalance / float64(remainingMonths)
				firstMonthInterest := baseBalance * monthlyRate
				loans[i].MonthlyPayment = principalPerMonth + firstMonthInterest

			case models.InterestFirst:
				// 先息后本
				loans[i].MonthlyPayment = baseBalance * monthlyRate
			}

			loans[i].MonthlyPayment = float64(int(loans[i].MonthlyPayment*100+0.5)) / 100
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

	// 打印从数据库读取的原始数据
	fmt.Printf("\n=== GetLoanByID Debug ===\n")
	fmt.Printf("Loan ID: %d\n", loan.ID)
	fmt.Printf("数据库中的月供: %.2f\n", loan.MonthlyPayment)
	fmt.Printf("数据库中的余额: %.2f\n", loan.RemainingAmount)
	fmt.Printf("还款记录数量: %d\n", len(loan.Payments))

	// 计算当前余额和月供
	now := time.Now()
	fmt.Printf("当前时间: %s\n", now.Format("2006-01-02"))

	// 找到最后一次实际还款日期
	var lastActualPaymentDate time.Time
	if len(loan.Payments) > 0 {
		for _, payment := range loan.Payments {
			if payment.PaymentDate.After(lastActualPaymentDate) {
				lastActualPaymentDate = payment.PaymentDate
			}
		}
		fmt.Printf("最后还款日期: %s\n", lastActualPaymentDate.Format("2006-01-02"))
	} else {
		fmt.Printf("没有还款记录\n")
	}

	// 余额：模拟到当前日期
	loan.RemainingAmount = s.calculatorService.CalculateRemainingAmount(loan, now)
	fmt.Printf("重新计算的余额: %.2f\n", loan.RemainingAmount)

	// 如果余额为0，更新状态为已完成
	if loan.RemainingAmount <= 0 && loan.Status != models.StatusCompleted {
		loan.Status = models.StatusCompleted
		loan.MonthlyPayment = 0
		fmt.Printf("余额为0，设置月供为0\n")
		s.db.Model(&loan).Updates(map[string]interface{}{
			"status":          models.StatusCompleted,
			"monthly_payment": 0,
		})
	} else {
		// 重新计算月供（基于最后还款后的余额，而不是当前余额）
		fmt.Printf("=== 开始重新计算月供 ===\n")

		// 计算最后还款后的余额（而不是当前余额）
		var baseBalance float64
		var baseDate time.Time
		if !lastActualPaymentDate.IsZero() {
			// 计算最后还款日的次日余额
			baseBalance = s.calculatorService.CalculateRemainingAmount(loan, lastActualPaymentDate.AddDate(0, 0, 1))
			baseDate = lastActualPaymentDate
			fmt.Printf("最后还款后余额: %.2f (还款日: %s)\n", baseBalance, lastActualPaymentDate.Format("2006-01-02"))
		} else {
			baseBalance = loan.TotalAmount
			baseDate = loan.StartDate
			fmt.Printf("无还款记录，使用初始金额: %.2f\n", baseBalance)
		}

		fmt.Printf("年利率: %.2f%%\n", loan.InterestRate)
		fmt.Printf("还款方式: %s\n", loan.PaymentMethod)

		// 计算下一个还款日
		var nextPaymentDate time.Time
		if loan.FirstPaymentDate != nil {
			nextPaymentDate = *loan.FirstPaymentDate
			fmt.Printf("使用第一期还款日: %s\n", nextPaymentDate.Format("2006-01-02"))
		} else {
			nextPaymentDate = loan.StartDate.AddDate(0, 1, 0)
			nextPaymentDate = time.Date(
				nextPaymentDate.Year(),
				nextPaymentDate.Month(),
				loan.PaymentDate,
				0, 0, 0, 0,
				nextPaymentDate.Location(),
			)
			fmt.Printf("计算的第一期还款日: %s\n", nextPaymentDate.Format("2006-01-02"))
		}

		// 找到baseDate之后的下一个还款日
		originalNext := nextPaymentDate
		for nextPaymentDate.Before(baseDate) || nextPaymentDate.Equal(baseDate) {
			nextPaymentDate = nextPaymentDate.AddDate(0, 1, 0)
		}
		fmt.Printf("从 %s 推进到 %s\n", originalNext.Format("2006-01-02"), nextPaymentDate.Format("2006-01-02"))

		// 计算最后还款日
		lastPaymentDate := time.Date(
			loan.EndDate.Year(),
			loan.EndDate.Month(),
			loan.PaymentDate,
			0, 0, 0, 0,
			loan.EndDate.Location(),
		)
		fmt.Printf("基于结束日期 %s 计算的最后还款日: %s\n",
			loan.EndDate.Format("2006-01-02"),
			lastPaymentDate.Format("2006-01-02"))

		if lastPaymentDate.After(loan.EndDate) {
			lastPaymentDate = lastPaymentDate.AddDate(0, -1, 0)
			fmt.Printf("调整后的最后还款日: %s\n", lastPaymentDate.Format("2006-01-02"))
		}

		// 计算剩余期数（需要+1因为包含首尾）
		remainingMonths := s.calculatorService.CalculateLoanMonths(nextPaymentDate, lastPaymentDate) + 1
		fmt.Printf("剩余期数: %d 个月\n", remainingMonths)

		// 计算月利率
		monthlyRate := s.calculatorService.CalculateMonthlyRate(loan.InterestRate)
		fmt.Printf("月利率: %.6f\n", monthlyRate)

		// 根据还款方式计算月供
		switch loan.PaymentMethod {
		case models.EqualInstallment:
			// 等额本息
			if monthlyRate == 0 {
				loan.MonthlyPayment = baseBalance / float64(remainingMonths)
			} else {
				pow := 1.0
				for i := 0; i < remainingMonths; i++ {
					pow *= (1 + monthlyRate)
				}
				loan.MonthlyPayment = baseBalance * monthlyRate * pow / (pow - 1)
			}

		case models.EqualPrincipal:
			// 等额本金：首月还款 = 每月本金 + 首月利息
			principalPerMonth := baseBalance / float64(remainingMonths)
			firstMonthInterest := baseBalance * monthlyRate
			loan.MonthlyPayment = principalPerMonth + firstMonthInterest

		case models.InterestFirst:
			// 先息后本：每月利息
			loan.MonthlyPayment = baseBalance * monthlyRate
		}

		// 保留两位小数
		loan.MonthlyPayment = float64(int(loan.MonthlyPayment*100+0.5)) / 100

		fmt.Printf("计算出的月供: %.2f\n", loan.MonthlyPayment)
		fmt.Printf("=== 月供计算完成 ===\n")
	}

	fmt.Printf("最终返回的月供: %.2f\n", loan.MonthlyPayment)
	fmt.Printf("=== End Debug ===\n\n")

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
	// 查询现有贷款以获取还款记录
	var existingLoan models.Loan
	if err := s.db.Preload("Payments").First(&existingLoan, loan.ID).Error; err != nil {
		return err
	}

	// 找到最后一次实际还款日期
	var lastActualPaymentDate time.Time
	if len(existingLoan.Payments) > 0 {
		for _, payment := range existingLoan.Payments {
			if payment.PaymentDate.After(lastActualPaymentDate) {
				lastActualPaymentDate = payment.PaymentDate
			}
		}
	}

	// 计算当前余额
	now := time.Now()
	existingLoan.TotalAmount = loan.TotalAmount
	existingLoan.InterestRate = loan.InterestRate
	existingLoan.StartDate = loan.StartDate
	existingLoan.EndDate = loan.EndDate
	existingLoan.PaymentDate = loan.PaymentDate
	existingLoan.PaymentMethod = loan.PaymentMethod
	existingLoan.Status = loan.Status
	existingLoan.BorrowerName = loan.BorrowerName
	existingLoan.BankName = loan.BankName
	if loan.FirstPaymentDate != nil {
		existingLoan.FirstPaymentDate = loan.FirstPaymentDate
	}

	// 重新计算余额
	existingLoan.RemainingAmount = s.calculatorService.CalculateRemainingAmount(existingLoan, now)

	// 如果余额为0，更新状态为已完成
	if existingLoan.RemainingAmount <= 0 {
		existingLoan.Status = models.StatusCompleted
		existingLoan.MonthlyPayment = 0
	} else {
		// 月供：基于最后实际还款后的余额计算
		var baseBalance float64
		var baseDate time.Time
		if !lastActualPaymentDate.IsZero() {
			baseBalance = s.calculatorService.CalculateRemainingAmount(existingLoan, lastActualPaymentDate.AddDate(0, 0, 1))
			baseDate = lastActualPaymentDate
		} else {
			baseBalance = existingLoan.TotalAmount
			baseDate = existingLoan.StartDate
		}

		// 计算下一个还款日
		var nextPaymentDate time.Time
		if existingLoan.FirstPaymentDate != nil {
			nextPaymentDate = *existingLoan.FirstPaymentDate
		} else {
			nextPaymentDate = existingLoan.StartDate.AddDate(0, 1, 0)
			nextPaymentDate = time.Date(
				nextPaymentDate.Year(),
				nextPaymentDate.Month(),
				existingLoan.PaymentDate,
				0, 0, 0, 0,
				nextPaymentDate.Location(),
			)
		}

		for nextPaymentDate.Before(baseDate) || nextPaymentDate.Equal(baseDate) {
			nextPaymentDate = nextPaymentDate.AddDate(0, 1, 0)
		}

		// 计算最后还款日
		lastPaymentDate := time.Date(
			existingLoan.EndDate.Year(),
			existingLoan.EndDate.Month(),
			existingLoan.PaymentDate,
			0, 0, 0, 0,
			existingLoan.EndDate.Location(),
		)
		if lastPaymentDate.After(existingLoan.EndDate) {
			lastPaymentDate = lastPaymentDate.AddDate(0, -1, 0)
		}

		existingLoan.MonthlyPayment = s.calculatorService.CalculateMonthlyPayment(
			baseBalance,
			existingLoan.InterestRate,
			nextPaymentDate,
			lastPaymentDate,
			existingLoan.PaymentMethod,
		)
	}

	return s.db.Save(&existingLoan).Error
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

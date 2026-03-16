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
			// 重新计算月供（基于当前余额和剩余期限）
			remainingMonths := s.calculatorService.CalculateLoanMonths(now, loans[i].EndDate)
			if remainingMonths > 0 {
				loans[i].MonthlyPayment = s.calculatorService.CalculateMonthlyPayment(
					loans[i].RemainingAmount,
					loans[i].InterestRate,
					now,
					loans[i].EndDate,
					loans[i].PaymentMethod,
				)
			}
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
		// 重新计算月供（基于当前余额和剩余期限）
		remainingMonths := s.calculatorService.CalculateLoanMonths(now, loan.EndDate)
		if remainingMonths > 0 {
			loan.MonthlyPayment = s.calculatorService.CalculateMonthlyPayment(
				loan.RemainingAmount,
				loan.InterestRate,
				now,
				loan.EndDate,
				loan.PaymentMethod,
			)
		}
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

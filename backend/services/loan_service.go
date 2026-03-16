package services

import (
	"loan-manager-wails/backend/models"

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
	err := s.db.Order("created_at DESC").Find(&loans).Error
	return loans, err
}

// GetLoanByID 根据ID获取贷款
func (s *LoanService) GetLoanByID(id uint) (*models.Loan, error) {
	var loan models.Loan
	err := s.db.Preload("Payments").First(&loan, id).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

// CreateLoan 创建贷款
func (s *LoanService) CreateLoan(loan *models.Loan) error {
	// 计算月供
	loan.MonthlyPayment = s.calculatorService.CalculateMonthlyPayment(
		loan.Amount,
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
		loan.Amount,
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

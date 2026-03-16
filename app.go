package main

import (
	"context"
	"loan-manager-wails/backend/database"
	"loan-manager-wails/backend/models"
	"loan-manager-wails/backend/services"
	"os"
	"path/filepath"
	"time"
)

// App struct
type App struct {
	ctx               context.Context
	loanService       *services.LoanService
	paymentService    *services.PaymentService
	calculatorService *services.CalculatorService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// 创建应用数据目录
	appDataDir := filepath.Join(homeDir, ".loan-manager")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		panic(err)
	}

	// 初始化数据库（存储在用户目录）
	dbPath := filepath.Join(appDataDir, "loans.db")
	if err := database.InitDB(dbPath); err != nil {
		panic(err)
	}

	// 初始化服务
	db := database.GetDB()
	a.loanService = services.NewLoanService(db)
	a.paymentService = services.NewPaymentService(db)
	a.calculatorService = services.NewCalculatorService()
}

// GetAllLoans 获取所有贷款
func (a *App) GetAllLoans() ([]models.Loan, error) {
	return a.loanService.GetAllLoans()
}

// GetLoanByID 根据ID获取贷款
func (a *App) GetLoanByID(id uint) (*models.Loan, error) {
	return a.loanService.GetLoanByID(id)
}

// CreateLoan 创建贷款
func (a *App) CreateLoan(loan models.Loan) error {
	return a.loanService.CreateLoan(&loan)
}

// UpdateLoan 更新贷款
func (a *App) UpdateLoan(loan models.Loan) error {
	return a.loanService.UpdateLoan(&loan)
}

// DeleteLoan 删除贷款
func (a *App) DeleteLoan(id uint) error {
	return a.loanService.DeleteLoan(id)
}

// MakePayment 记录还款
func (a *App) MakePayment(loanID uint, amount float64, paymentDateStr string) error {
	// 解析还款时间
	paymentDate, err := time.Parse(time.RFC3339, paymentDateStr)
	if err != nil {
		return err
	}
	return a.paymentService.MakePayment(loanID, amount, paymentDate)
}

// GetPaymentSchedule 获取还款计划表
func (a *App) GetPaymentSchedule(loanID uint) ([]services.PaymentDetail, error) {
	// 获取贷款信息
	loan, err := a.loanService.GetLoanByID(loanID)
	if err != nil {
		return nil, err
	}

	// 生成还款计划
	now := time.Now()
	schedule := a.calculatorService.GeneratePaymentSchedule(*loan, now)
	return schedule, nil
}

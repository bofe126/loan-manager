package models

// PaymentMethod 还款方式
type PaymentMethod string

const (
	EqualInstallment PaymentMethod = "等额本息" // 等额本息
	EqualPrincipal   PaymentMethod = "等额本金" // 等额本金
	InterestFirst    PaymentMethod = "先息后本" // 先息后本
)

// LoanStatus 贷款状态
type LoanStatus string

const (
	StatusActive    LoanStatus = "active"    // 进行中
	StatusCompleted LoanStatus = "completed" // 已完成
	StatusOverdue   LoanStatus = "overdue"   // 逾期
)

// 计算常量
const (
	MonthsPerYear  = 12
	Percentage     = 100
	DecimalPlaces  = 2
)

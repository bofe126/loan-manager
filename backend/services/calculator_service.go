package services

import (
	"loan-manager-wails/backend/models"
	"math"
	"time"
)

// CalculatorService 计算服务
type CalculatorService struct{}

// NewCalculatorService 创建计算服务
func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

// CalculateLoanMonths 计算贷款期数（月）
func (s *CalculatorService) CalculateLoanMonths(startDate, endDate time.Time) int {
	years := endDate.Year() - startDate.Year()
	months := int(endDate.Month() - startDate.Month())
	totalMonths := years*models.MonthsPerYear + months

	// 如果结束日期的天数小于开始日期，减去一个月
	if endDate.Day() < startDate.Day() {
		totalMonths--
	}

	if totalMonths < 1 {
		totalMonths = 1
	}

	return totalMonths
}

// CalculateMonthlyRate 计算月利率
func (s *CalculatorService) CalculateMonthlyRate(annualRate float64) float64 {
	return annualRate / models.Percentage / models.MonthsPerYear
}

// CalculateMonthlyPayment 计算月供
func (s *CalculatorService) CalculateMonthlyPayment(amount, interestRate float64, startDate, endDate time.Time, paymentMethod models.PaymentMethod) float64 {
	months := s.CalculateLoanMonths(startDate, endDate)
	monthlyRate := s.CalculateMonthlyRate(interestRate)

	var monthlyPayment float64

	switch paymentMethod {
	case models.EqualInstallment:
		// 等额本息：PMT = P * r * (1+r)^n / ((1+r)^n - 1)
		if monthlyRate == 0 {
			monthlyPayment = amount / float64(months)
		} else {
			pow := math.Pow(1+monthlyRate, float64(months))
			monthlyPayment = amount * monthlyRate * pow / (pow - 1)
		}

	case models.EqualPrincipal:
		// 等额本金：首月还款 = 每月本金 + 首月利息
		principalPerMonth := amount / float64(months)
		firstMonthInterest := amount * monthlyRate
		monthlyPayment = principalPerMonth + firstMonthInterest

	case models.InterestFirst:
		// 先息后本：每月利息
		monthlyPayment = amount * monthlyRate
	}

	// 保留两位小数
	return math.Round(monthlyPayment*100) / 100
}

// CalculateRemainingAmount 计算剩余金额（基于当前日期）
func (s *CalculatorService) CalculateRemainingAmount(loan models.Loan, currentDate time.Time) float64 {
	// 计算已过月数
	monthsPassed := s.CalculateLoanMonths(loan.StartDate, currentDate)
	totalMonths := s.CalculateLoanMonths(loan.StartDate, loan.EndDate)

	if monthsPassed >= totalMonths {
		return 0
	}

	if monthsPassed < 0 {
		monthsPassed = 0
	}

	monthlyRate := s.CalculateMonthlyRate(loan.InterestRate)
	var remaining float64

	switch loan.PaymentMethod {
	case models.EqualInstallment:
		// 等额本息：剩余本金 = P * ((1+r)^n - (1+r)^m) / ((1+r)^n - 1)
		if monthlyRate == 0 {
			remaining = loan.Amount * float64(totalMonths-monthsPassed) / float64(totalMonths)
		} else {
			powTotal := math.Pow(1+monthlyRate, float64(totalMonths))
			powPassed := math.Pow(1+monthlyRate, float64(monthsPassed))
			remaining = loan.Amount * (powTotal - powPassed) / (powTotal - 1)
		}

	case models.EqualPrincipal:
		// 等额本金：剩余本金 = P - (P / n) * m
		principalPerMonth := loan.Amount / float64(totalMonths)
		remaining = loan.Amount - principalPerMonth*float64(monthsPassed)

	case models.InterestFirst:
		// 先息后本：剩余本金 = P（本金不变）
		remaining = loan.Amount
	}

	if remaining < 0 {
		remaining = 0
	}

	// 保留两位小数
	return math.Round(remaining*100) / 100
}

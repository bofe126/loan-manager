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

// CalculateRemainingAmount 计算剩余金额
// 计算逻辑：逐月模拟还款过程，考虑正常还款和额外还款的影响
// 提前还款时先还利息，剩余部分还本金；不足月按天计算利息
func (s *CalculatorService) CalculateRemainingAmount(loan models.Loan, currentDate time.Time) float64 {
	totalMonths := s.CalculateLoanMonths(loan.StartDate, loan.EndDate)

	// 确定第一期还款日期
	var firstPaymentDate time.Time
	if loan.FirstPaymentDate != nil {
		firstPaymentDate = *loan.FirstPaymentDate
	} else {
		firstPaymentDate = loan.StartDate.AddDate(0, 1, 0)
		firstPaymentDate = time.Date(
			firstPaymentDate.Year(),
			firstPaymentDate.Month(),
			loan.PaymentDate,
			0, 0, 0, 0,
			firstPaymentDate.Location(),
		)
	}

	monthsPassed := s.CalculateLoanMonths(firstPaymentDate, currentDate) + 1

	if monthsPassed >= totalMonths {
		return 0
	}

	if monthsPassed <= 0 {
		// 还没开始还款，但可能有提前还款
		remaining := loan.TotalAmount
		for _, payment := range loan.Payments {
			days := int(payment.PaymentDate.Sub(loan.StartDate).Hours() / 24)
			dailyRate := s.CalculateMonthlyRate(loan.InterestRate) / 30
			interest := remaining * dailyRate * float64(days)

			principalPayment := payment.Amount - interest
			if principalPayment > 0 {
				remaining -= principalPayment
			}
		}
		if remaining < 0 {
			remaining = 0
		}
		return math.Round(remaining*100) / 100
	}

	monthlyRate := s.CalculateMonthlyRate(loan.InterestRate)
	dailyRate := monthlyRate / 30
	currentPrincipal := loan.TotalAmount
	remainingMonths := totalMonths
	lastPaymentDate := loan.StartDate

	// 逐月模拟还款过程
	for month := 1; month <= monthsPassed; month++ {
		if currentPrincipal <= 0 {
			break
		}

		// 计算本月还款日期
		var monthlyPaymentDate time.Time
		if month == 1 {
			monthlyPaymentDate = firstPaymentDate
		} else {
			monthlyPaymentDate = firstPaymentDate.AddDate(0, month-1, 0)
		}

		// 计算本月正常还款的本金部分
		var principalPayment float64
		days := int(monthlyPaymentDate.Sub(lastPaymentDate).Hours() / 24)

		switch loan.PaymentMethod {
		case models.EqualInstallment:
			if monthlyRate == 0 {
				principalPayment = currentPrincipal / float64(remainingMonths)
			} else {
				currentMonthlyPayment := s.calculateMonthlyPaymentForAmount(
					currentPrincipal,
					monthlyRate,
					remainingMonths,
				)
				// 如果是第一期且天数不是标准30天，按实际天数计算利息
				var monthlyInterest float64
				if month == 1 && days != 30 {
					monthlyInterest = currentPrincipal * dailyRate * float64(days)
				} else {
					monthlyInterest = currentPrincipal * monthlyRate
				}
				principalPayment = currentMonthlyPayment - monthlyInterest
			}

		case models.EqualPrincipal:
			principalPayment = loan.TotalAmount / float64(totalMonths)

		case models.InterestFirst:
			if month == totalMonths {
				principalPayment = currentPrincipal
			} else {
				principalPayment = 0
			}
		}

		currentPrincipal -= principalPayment
		if currentPrincipal < 0 {
			currentPrincipal = 0
		}
		lastPaymentDate = monthlyPaymentDate

		// 检查本月还款日之后是否有额外还款
		var nextMonthPaymentDate time.Time
		if month == monthsPassed {
			nextMonthPaymentDate = currentDate.AddDate(0, 1, 0)
		} else {
			nextMonthPaymentDate = firstPaymentDate.AddDate(0, month, 0)
		}

		for _, payment := range loan.Payments {
			if payment.PaymentDate.After(monthlyPaymentDate) &&
			   (payment.PaymentDate.Before(nextMonthPaymentDate) || payment.PaymentDate.Equal(nextMonthPaymentDate)) {

				days := int(payment.PaymentDate.Sub(lastPaymentDate).Hours() / 24)
				interest := currentPrincipal * dailyRate * float64(days)

				principalPayment := payment.Amount - interest
				if principalPayment > 0 {
					currentPrincipal -= principalPayment
					if currentPrincipal < 0 {
						currentPrincipal = 0
					}
				}

				lastPaymentDate = payment.PaymentDate

				if currentPrincipal <= 0 {
					break
				}
			}
		}

		if currentPrincipal <= 0 {
			break
		}

		remainingMonths--
	}

	return math.Round(currentPrincipal*100) / 100
}

// calculateMonthlyPaymentForAmount 计算指定本金和期限的月供（等额本息）
func (s *CalculatorService) calculateMonthlyPaymentForAmount(principal float64, monthlyRate float64, months int) float64 {
	if monthlyRate == 0 {
		return principal / float64(months)
	}
	pow := math.Pow(1+monthlyRate, float64(months))
	return principal * monthlyRate * pow / (pow - 1)
}

// PaymentDetail 还款明细
type PaymentDetail struct {
	Month            int       `json:"month"`              // 第几期
	PaymentDate      time.Time `json:"payment_date"`       // 还款日期
	Principal        float64   `json:"principal"`          // 本金
	Interest         float64   `json:"interest"`           // 利息
	TotalPayment     float64   `json:"total_payment"`      // 总还款
	RemainingBalance float64   `json:"remaining_balance"`  // 剩余本金
	ExtraPayment     float64   `json:"extra_payment"`      // 额外还款
	ExtraInterest    float64   `json:"extra_interest"`     // 额外还款利息
	Note             string    `json:"note"`               // 备注
}

// GeneratePaymentSchedule 生成还款计划表（从开始到当前日期）
// 包含正常还款和提前还款，每笔还款都是单独一行
func (s *CalculatorService) GeneratePaymentSchedule(loan models.Loan, currentDate time.Time) []PaymentDetail {
	var schedule []PaymentDetail

	totalMonths := s.CalculateLoanMonths(loan.StartDate, loan.EndDate)

	// 确定第一期还款日期
	var firstPaymentDate time.Time
	if loan.FirstPaymentDate != nil {
		firstPaymentDate = *loan.FirstPaymentDate
	} else {
		firstPaymentDate = loan.StartDate.AddDate(0, 1, 0)
		firstPaymentDate = time.Date(
			firstPaymentDate.Year(),
			firstPaymentDate.Month(),
			loan.PaymentDate,
			0, 0, 0, 0,
			firstPaymentDate.Location(),
		)
	}

	monthsPassed := s.CalculateLoanMonths(firstPaymentDate, currentDate) + 1

	if monthsPassed <= 0 {
		return schedule
	}

	if monthsPassed > totalMonths {
		monthsPassed = totalMonths
	}

	monthlyRate := s.CalculateMonthlyRate(loan.InterestRate)
	dailyRate := monthlyRate / 30
	currentPrincipal := loan.TotalAmount
	remainingMonths := totalMonths
	lastPaymentDate := loan.StartDate
	periodNumber := 0

	// 逐月生成还款明细
	for month := 1; month <= monthsPassed; month++ {
		if currentPrincipal <= 0 {
			break
		}

		// 计算本月还款日期
		var monthlyPaymentDate time.Time
		if month == 1 {
			monthlyPaymentDate = firstPaymentDate
		} else {
			monthlyPaymentDate = firstPaymentDate.AddDate(0, month-1, 0)
		}

		// 计算本月正常还款
		var principalPayment float64
		var interestPayment float64
		days := int(monthlyPaymentDate.Sub(lastPaymentDate).Hours() / 24)

		switch loan.PaymentMethod {
		case models.EqualInstallment:
			if monthlyRate == 0 {
				principalPayment = currentPrincipal / float64(remainingMonths)
				interestPayment = 0
			} else {
				if month == 1 && days != 30 {
					interestPayment = currentPrincipal * dailyRate * float64(days)
					currentMonthlyPayment := s.calculateMonthlyPaymentForAmount(
						currentPrincipal,
						monthlyRate,
						remainingMonths,
					)
					principalPayment = currentMonthlyPayment - (currentPrincipal * monthlyRate)
				} else {
					currentMonthlyPayment := s.calculateMonthlyPaymentForAmount(
						currentPrincipal,
						monthlyRate,
						remainingMonths,
					)
					interestPayment = currentPrincipal * monthlyRate
					principalPayment = currentMonthlyPayment - interestPayment
				}
			}

		case models.EqualPrincipal:
			principalPayment = loan.TotalAmount / float64(totalMonths)
			if month == 1 && days != 30 {
				interestPayment = currentPrincipal * dailyRate * float64(days)
			} else {
				interestPayment = currentPrincipal * monthlyRate
			}

		case models.InterestFirst:
			if month == totalMonths {
				principalPayment = currentPrincipal
				interestPayment = 0
			} else {
				principalPayment = 0
				if month == 1 && days != 30 {
					interestPayment = currentPrincipal * dailyRate * float64(days)
				} else {
					interestPayment = currentPrincipal * monthlyRate
				}
			}
		}

		periodNumber++
		detail := PaymentDetail{
			Month:            periodNumber,
			PaymentDate:      monthlyPaymentDate,
			Principal:        math.Round(principalPayment*100) / 100,
			Interest:         math.Round(interestPayment*100) / 100,
			TotalPayment:     math.Round((principalPayment+interestPayment)*100) / 100,
			RemainingBalance: 0,
			Note:             "",
		}

		currentPrincipal -= principalPayment
		if currentPrincipal < 0 {
			currentPrincipal = 0
		}
		detail.RemainingBalance = math.Round(currentPrincipal*100) / 100
		lastPaymentDate = monthlyPaymentDate

		schedule = append(schedule, detail)

		// 检查本月还款日之后是否有额外还款
		var nextMonthPaymentDate time.Time
		if month == monthsPassed {
			nextMonthPaymentDate = currentDate.AddDate(0, 1, 0)
		} else {
			nextMonthPaymentDate = firstPaymentDate.AddDate(0, month, 0)
		}

		for _, payment := range loan.Payments {
			if payment.PaymentDate.After(monthlyPaymentDate) &&
			   (payment.PaymentDate.Before(nextMonthPaymentDate) || payment.PaymentDate.Equal(nextMonthPaymentDate)) {

				// 计算额外还款的利息（按天）
				days := int(payment.PaymentDate.Sub(lastPaymentDate).Hours() / 24)
				extraInterest := currentPrincipal * dailyRate * float64(days)
				extraInterest = math.Round(extraInterest*100) / 100

				// 本金还款
				extraPrincipal := payment.Amount - extraInterest
				if extraPrincipal < 0 {
					extraPrincipal = 0
				}

				periodNumber++
				extraDetail := PaymentDetail{
					Month:            periodNumber,
					PaymentDate:      payment.PaymentDate,
					Principal:        math.Round(extraPrincipal*100) / 100,
					Interest:         extraInterest,
					TotalPayment:     payment.Amount,
					RemainingBalance: 0,
					Note:             "提前还款",
				}

				if extraPrincipal > 0 {
					currentPrincipal -= extraPrincipal
					if currentPrincipal < 0 {
						currentPrincipal = 0
					}
				}
				extraDetail.RemainingBalance = math.Round(currentPrincipal*100) / 100

				lastPaymentDate = payment.PaymentDate

				schedule = append(schedule, extraDetail)

				if currentPrincipal <= 0 {
					break
				}
			}
		}

		if currentPrincipal <= 0 {
			break
		}

		remainingMonths--
	}

	return schedule
}

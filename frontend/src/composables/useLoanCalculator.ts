import { PaymentMethod } from '../types/loan';

const MONTHS_PER_YEAR = 12;
const PERCENTAGE = 100;

export function useLoanCalculator() {
  // 计算贷款期数（月）
  const calculateLoanMonths = (startDate: Date, endDate: Date): number => {
    const years = endDate.getFullYear() - startDate.getFullYear();
    const months = endDate.getMonth() - startDate.getMonth();
    let totalMonths = years * MONTHS_PER_YEAR + months;

    if (endDate.getDate() < startDate.getDate()) {
      totalMonths--;
    }

    return totalMonths < 1 ? 1 : totalMonths;
  };

  // 计算月利率
  const calculateMonthlyRate = (annualRate: number): number => {
    return annualRate / PERCENTAGE / MONTHS_PER_YEAR;
  };

  // 计算月供
  const calculateMonthlyPayment = (
    amount: number,
    interestRate: number,
    startDate: Date,
    endDate: Date,
    paymentMethod: PaymentMethod
  ): number => {
    const months = calculateLoanMonths(startDate, endDate);
    const monthlyRate = calculateMonthlyRate(interestRate);

    let monthlyPayment = 0;

    switch (paymentMethod) {
      case '等额本息': {
        // 等额本息：PMT = P * r * (1+r)^n / ((1+r)^n - 1)
        if (monthlyRate === 0) {
          monthlyPayment = amount / months;
        } else {
          const pow = Math.pow(1 + monthlyRate, months);
          monthlyPayment = (amount * monthlyRate * pow) / (pow - 1);
        }
        break;
      }
      case '等额本金': {
        // 等额本金：首月还款 = 每月本金 + 首月利息
        const principalPerMonth = amount / months;
        const firstMonthInterest = amount * monthlyRate;
        monthlyPayment = principalPerMonth + firstMonthInterest;
        break;
      }
      case '先息后本': {
        // 先息后本：每月利息
        monthlyPayment = amount * monthlyRate;
        break;
      }
    }

    // 保留两位小数
    return Math.round(monthlyPayment * 100) / 100;
  };

  // 计算剩余金额
  const calculateRemainingAmount = (
    amount: number,
    interestRate: number,
    startDate: Date,
    endDate: Date,
    paymentMethod: PaymentMethod,
    currentDate: Date = new Date()
  ): number => {
    const monthsPassed = calculateLoanMonths(startDate, currentDate);
    const totalMonths = calculateLoanMonths(startDate, endDate);

    if (monthsPassed >= totalMonths) {
      return 0;
    }

    if (monthsPassed < 0) {
      return amount;
    }

    const monthlyRate = calculateMonthlyRate(interestRate);
    let remaining = 0;

    switch (paymentMethod) {
      case '等额本息': {
        // 等额本息：剩余本金 = P * ((1+r)^n - (1+r)^m) / ((1+r)^n - 1)
        if (monthlyRate === 0) {
          remaining = (amount * (totalMonths - monthsPassed)) / totalMonths;
        } else {
          const powTotal = Math.pow(1 + monthlyRate, totalMonths);
          const powPassed = Math.pow(1 + monthlyRate, monthsPassed);
          remaining = (amount * (powTotal - powPassed)) / (powTotal - 1);
        }
        break;
      }
      case '等额本金': {
        // 等额本金：剩余本金 = P - (P / n) * m
        const principalPerMonth = amount / totalMonths;
        remaining = amount - principalPerMonth * monthsPassed;
        break;
      }
      case '先息后本': {
        // 先息后本：剩余本金 = P（本金不变）
        remaining = amount;
        break;
      }
    }

    if (remaining < 0) {
      remaining = 0;
    }

    // 保留两位小数
    return Math.round(remaining * 100) / 100;
  };

  return {
    calculateLoanMonths,
    calculateMonthlyRate,
    calculateMonthlyPayment,
    calculateRemainingAmount,
  };
}

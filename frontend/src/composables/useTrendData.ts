import { Loan, PaymentMethod } from '../types/loan';
import { useLoanCalculator } from './useLoanCalculator';

export function useTrendData() {
  const { calculateLoanMonths, calculateMonthlyRate } = useLoanCalculator();

  // 计算未来24个月的趋势数据
  const calculateTrendData = (loans: Loan[]) => {
    const months = 24;
    const labels: string[] = [];
    const monthlyPayments: number[] = [];
    const remainingAmounts: number[] = [];

    const now = new Date();

    for (let i = 0; i < months; i++) {
      const futureDate = new Date(now.getFullYear(), now.getMonth() + i, 1);
      labels.push(`${futureDate.getFullYear()}-${String(futureDate.getMonth() + 1).padStart(2, '0')}`);

      let totalMonthlyPayment = 0;
      let totalRemaining = 0;

      loans.forEach((loan) => {
        const startDate = new Date(loan.start_date);
        const endDate = new Date(loan.end_date);

        // 检查贷款是否在此月份内有效
        if (futureDate >= startDate && futureDate <= endDate) {
          totalMonthlyPayment += loan.monthly_payment;

          // 计算该月的剩余金额
          const monthsPassed = calculateLoanMonths(startDate, futureDate);
          const totalMonths = calculateLoanMonths(startDate, endDate);
          const monthlyRate = calculateMonthlyRate(loan.interest_rate);

          let remaining = 0;

          switch (loan.payment_method as PaymentMethod) {
            case '等额本息': {
              if (monthlyRate === 0) {
                remaining = (loan.amount * (totalMonths - monthsPassed)) / totalMonths;
              } else {
                const powTotal = Math.pow(1 + monthlyRate, totalMonths);
                const powPassed = Math.pow(1 + monthlyRate, monthsPassed);
                remaining = (loan.amount * (powTotal - powPassed)) / (powTotal - 1);
              }
              break;
            }
            case '等额本金': {
              const principalPerMonth = loan.amount / totalMonths;
              remaining = loan.amount - principalPerMonth * monthsPassed;
              break;
            }
            case '先息后本': {
              remaining = loan.amount;
              break;
            }
          }

          if (remaining < 0) remaining = 0;
          totalRemaining += Math.round(remaining * 100) / 100;
        }
      });

      monthlyPayments.push(Math.round(totalMonthlyPayment * 100) / 100);
      remainingAmounts.push(Math.round(totalRemaining * 100) / 100);
    }

    return {
      labels,
      monthlyPayments,
      remainingAmounts,
    };
  };

  return {
    calculateTrendData,
  };
}

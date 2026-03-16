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
        if (loan.status !== 'active') return;

        const startDate = new Date(loan.start_date);
        const endDate = new Date(loan.end_date);

        // 检查贷款是否在此月份内有效
        if (futureDate >= now && futureDate <= endDate) {
          // 使用当前的月供（已经考虑了提前还款）
          totalMonthlyPayment += loan.monthly_payment;

          // 计算该月的剩余金额（从当前余额开始计算）
          const monthsFromNow = calculateLoanMonths(now, futureDate);
          const remainingMonthsTotal = calculateLoanMonths(now, endDate);
          const monthlyRate = calculateMonthlyRate(loan.interest_rate);

          // 使用当前剩余金额作为起点
          const currentRemaining = loan.remaining_amount || loan.total_amount;
          let remaining = 0;

          if (monthsFromNow === 0) {
            // 当前月份，使用当前余额
            remaining = currentRemaining;
          } else {
            // 未来月份，基于当前余额计算
            switch (loan.payment_method as PaymentMethod) {
              case '等额本息': {
                if (monthlyRate === 0) {
                  remaining = (currentRemaining * (remainingMonthsTotal - monthsFromNow)) / remainingMonthsTotal;
                } else {
                  const powTotal = Math.pow(1 + monthlyRate, remainingMonthsTotal);
                  const powPassed = Math.pow(1 + monthlyRate, monthsFromNow);
                  remaining = (currentRemaining * (powTotal - powPassed)) / (powTotal - 1);
                }
                break;
              }
              case '等额本金': {
                const principalPerMonth = currentRemaining / remainingMonthsTotal;
                remaining = currentRemaining - principalPerMonth * monthsFromNow;
                break;
              }
              case '先息后本': {
                // 先息后本：最后一个月才还本金
                if (monthsFromNow >= remainingMonthsTotal) {
                  remaining = 0;
                } else {
                  remaining = currentRemaining;
                }
                break;
              }
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

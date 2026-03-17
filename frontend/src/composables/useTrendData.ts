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
    // 获取当前月份的第一天，用于比较
    const currentMonthStart = new Date(now.getFullYear(), now.getMonth(), 1);

    for (let i = 0; i < months; i++) {
      const futureDate = new Date(now.getFullYear(), now.getMonth() + i, 1);
      labels.push(`${futureDate.getFullYear()}-${String(futureDate.getMonth() + 1).padStart(2, '0')}`);

      let totalMonthlyPayment = 0;
      let totalRemaining = 0;

      loans.forEach((loan) => {
        if (loan.status !== 'active') return;

        const startDate = new Date(loan.start_date);
        const endDate = new Date(loan.end_date);
        const endMonthStart = new Date(endDate.getFullYear(), endDate.getMonth(), 1);

        // 检查贷款是否在此月份内有效
        if (futureDate >= currentMonthStart && futureDate <= endMonthStart) {
          // 使用当前的月供（已经考虑了提前还款）
          totalMonthlyPayment += loan.monthly_payment;

          // 计算该月的剩余金额
          // 对于当前月份，需要反推到月初的余额
          const currentRemaining = loan.remaining_amount || loan.total_amount;
          let remaining = 0;

          if (i === 0) {
            // 当前月份：从当前余额反推到月初
            // 如果本月已经有还款，需要加回本金部分
            const currentDay = now.getDate();
            const paymentDay = loan.payment_date || 1;

            if (currentDay > paymentDay) {
              // 本月还款日已过，加回本月还款的本金部分
              const monthlyRate = calculateMonthlyRate(loan.interest_rate);
              const interest = currentRemaining * monthlyRate / (1 - monthlyRate);
              const principal = loan.monthly_payment - interest;
              remaining = currentRemaining + principal;
            } else {
              // 本月还款日未到，当前余额就是月初余额
              remaining = currentRemaining;
            }
          } else {
            // 未来月份：从当前余额投影
            const monthsFromNow = i;
            const remainingMonthsTotal = calculateLoanMonths(now, endDate);
            const monthlyRate = calculateMonthlyRate(loan.interest_rate);

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

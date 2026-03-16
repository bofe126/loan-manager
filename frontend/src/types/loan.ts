export interface Loan {
  id: number;
  borrower_name: string;
  bank_name: string;
  total_amount: number;
  remaining_amount: number;
  interest_rate: number;
  start_date: string;
  end_date: string;
  first_payment_date?: string;
  payment_date: number;
  status: LoanStatus;
  payment_method: PaymentMethod;
  monthly_payment: number;
  created_at: string;
  payments?: Payment[];
}

export interface Payment {
  id: number;
  loan_id: number;
  amount: number;
  payment_date: string;
}

export type PaymentMethod = '等额本息' | '等额本金' | '先息后本';

export type LoanStatus = 'active' | 'completed' | 'overdue';

export const PaymentMethods = {
  EQUAL_INSTALLMENT: '等额本息' as PaymentMethod,
  EQUAL_PRINCIPAL: '等额本金' as PaymentMethod,
  INTEREST_FIRST: '先息后本' as PaymentMethod,
};

export const LoanStatuses = {
  ACTIVE: 'active' as LoanStatus,
  COMPLETED: 'completed' as LoanStatus,
  OVERDUE: 'overdue' as LoanStatus,
};

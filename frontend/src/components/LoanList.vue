<template>
  <div class="loan-list-container">
    <div v-if="loading" class="text-center loading-state">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">加载中...</span>
      </div>
    </div>
    <div v-else-if="loans.length === 0" class="empty-state">
      <p>暂无贷款记录</p>
    </div>
    <div v-else class="table-wrapper">
      <table class="loan-table">
        <thead>
          <tr>
            <th class="text-left">银行</th>
            <th class="text-center">贷款余额</th>
            <th class="text-center">利率</th>
            <th class="text-center">月供</th>
            <th class="text-center">还款方式</th>
            <th class="text-center">开始日期</th>
            <th class="text-center">结束日期</th>
            <th class="text-center">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="loan in loans" :key="loan.id" class="loan-row">
            <td class="text-left">
              <span class="bank-name">{{ loan.bank_name }}</span>
            </td>
            <td class="text-center font-mono">¥{{ Math.round(loan.remaining_amount).toLocaleString('zh-CN') }}</td>
            <td class="text-center font-mono">{{ loan.interest_rate.toFixed(2) }}%</td>
            <td class="text-center font-mono primary-text">¥{{ Math.round(loan.monthly_payment).toLocaleString('zh-CN') }}</td>
            <td class="text-center">{{ getPaymentMethodText(loan.payment_method) }}</td>
            <td class="text-center font-mono text-muted">{{ formatDate(loan.start_date) }}</td>
            <td class="text-center font-mono text-muted">{{ formatDate(loan.end_date) }}</td>
            <td class="text-center">
              <div class="action-buttons">
                <button class="btn-icon btn-schedule" @click="$emit('schedule', loan.id)" title="还款明细">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                    <polyline points="14 2 14 8 20 8"/>
                    <line x1="16" y1="13" x2="8" y2="13"/>
                    <line x1="16" y1="17" x2="8" y2="17"/>
                    <polyline points="10 9 9 9 8 9"/>
                  </svg>
                </button>
                <button class="btn-icon btn-edit" @click="$emit('edit', loan.id)" title="编辑">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                  </svg>
                </button>
                <button class="btn-icon btn-payment" @click="$emit('payment', loan.id)" :disabled="loan.status === 'completed'" title="还款">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
                  </svg>
                </button>
                <button class="btn-icon btn-delete" @click="handleDelete(loan.id)" title="删除">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { models } from '../../wailsjs/go/models';

interface Props {
  loans: models.Loan[];
  loading?: boolean;
}

interface Emits {
  (e: 'edit', id: number): void;
  (e: 'delete', id: number): void;
  (e: 'payment', id: number): void;
  (e: 'schedule', id: number): void;
}

defineProps<Props>();
const emit = defineEmits<Emits>();

const formatDate = (dateStr: any): string => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN');
};

const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    active: '进行中',
    completed: '已完成',
    overdue: '逾期',
  };
  return statusMap[status] || status;
};

const getStatusClass = (status: string): string => {
  const classMap: Record<string, string> = {
    active: 'status-badge status-active',
    completed: 'status-badge status-completed',
    overdue: 'status-badge status-overdue',
  };
  return classMap[status] || 'status-badge';
};

const getPaymentMethodText = (method: string): string => {
  const methodMap: Record<string, string> = {
    equal_principal_interest: '等额本息',
    equal_principal: '等额本金',
    interest_first: '先息后本',
  };
  return methodMap[method] || method;
};

const handleDelete = (id: number) => {
  if (confirm('确定要删除这条贷款记录吗？')) {
    emit('delete', id);
  }
};
</script>

<style scoped>
.loan-list-container {
  width: 100%;
}

.loading-state,
.empty-state {
  padding: 2rem;
  text-align: center;
  color: var(--text-secondary);
}

.table-wrapper {
  background: var(--glass-bg);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-lg);
  overflow-y: auto;
  max-height: calc(100vh - 8rem);
}

.loan-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8125rem;
  table-layout: fixed;
}

.loan-table th:nth-child(1),
.loan-table td:nth-child(1) { width: 15%; }  /* 银行 */

.loan-table th:nth-child(2),
.loan-table td:nth-child(2) { width: 13%; }  /* 贷款余额 */

.loan-table th:nth-child(3),
.loan-table td:nth-child(3) { width: 8%; }   /* 利率 */

.loan-table th:nth-child(4),
.loan-table td:nth-child(4) { width: 13%; }  /* 月供 */

.loan-table th:nth-child(5),
.loan-table td:nth-child(5) { width: 10%; }  /* 还款方式 */

.loan-table th:nth-child(6),
.loan-table td:nth-child(6) { width: 10%; }  /* 开始日期 */

.loan-table th:nth-child(7),
.loan-table td:nth-child(7) { width: 10%; }  /* 结束日期 */

.loan-table th:nth-child(8),
.loan-table td:nth-child(8) { width: 21%; }  /* 操作 */

.loan-table thead {
  background: rgba(255, 255, 255, 0.08);
  position: sticky;
  top: 0;
  z-index: 10;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(12px);
}

.loan-table th {
  padding: 0.625rem 0.75rem;
  font-weight: 600;
  font-size: 0.8125rem;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  white-space: nowrap;
  border-bottom: 1px solid var(--glass-border);
  line-height: 1.2;
}

.loan-table th.text-left {
  text-align: left;
}

.loan-table th.text-center {
  text-align: center;
}

.loan-table tbody tr {
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  transition: all var(--transition-base);
}

.loan-table tbody tr:last-child {
  border-bottom: none;
}

.loan-table tbody tr:hover {
  background: rgba(255, 255, 255, 0.03);
}

.loan-table td {
  padding: 0.625rem 0.75rem;
  color: var(--text-primary);
  white-space: nowrap;
  line-height: 1.2;
}

.loan-table td.text-left {
  text-align: left;
}

.loan-table td.text-center {
  text-align: center;
}

.bank-name {
  font-weight: 600;
  color: var(--text-primary);
}

.font-mono {
  font-family: var(--font-mono);
}

.text-muted {
  color: var(--text-secondary);
}

.primary-text {
  color: #0080ff;
  font-weight: 600;
}

.status-badge {
  display: inline-block;
  padding: 0.3125rem 0.75rem;
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.status-active {
  background: rgba(57, 255, 20, 0.15);
  color: #39ff14;
  border: 1px solid rgba(57, 255, 20, 0.3);
}

.status-completed {
  background: rgba(100, 116, 139, 0.15);
  color: #94a3b8;
  border: 1px solid rgba(100, 116, 139, 0.3);
}

.status-overdue {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.action-buttons {
  display: flex;
  gap: 0.375rem;
  justify-content: center;
}

.btn-icon {
  padding: 0.375rem;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-base);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: transparent;
}

.btn-icon svg {
  display: block;
}

.btn-schedule {
  color: #a78bfa;
}

.btn-schedule:hover {
  background: rgba(167, 139, 250, 0.15);
  transform: scale(1.1);
}

.btn-edit {
  color: #0080ff;
}

.btn-edit:hover {
  background: rgba(0, 128, 255, 0.15);
  transform: scale(1.1);
}

.btn-payment {
  color: #39ff14;
}

.btn-payment:hover:not(:disabled) {
  background: rgba(57, 255, 20, 0.15);
  transform: scale(1.1);
}

.btn-payment:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.btn-delete {
  color: #ef4444;
}

.btn-delete:hover {
  background: rgba(239, 68, 68, 0.15);
  transform: scale(1.1);
}

@media (max-width: 1200px) {
  .table-wrapper {
    overflow-x: auto;
  }

  .loan-table {
    min-width: 900px;
  }
}
</style>

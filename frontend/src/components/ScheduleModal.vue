<template>
  <div v-if="show" class="modal-overlay" @click="handleClose">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h5 class="modal-title">还款明细</h5>
        <button class="btn-close" @click="handleClose">&times;</button>
      </div>
      <div class="modal-body">
        <div v-if="loading" class="text-center">
          <div class="spinner-border" role="status">
            <span class="visually-hidden">加载中...</span>
          </div>
        </div>
        <div v-else-if="schedule.length === 0" class="empty-state">
          暂无还款记录
        </div>
        <div v-else class="schedule-table-wrapper">
          <table class="schedule-table">
            <thead>
              <tr>
                <th>期数</th>
                <th>还款日期</th>
                <th class="text-right">本金</th>
                <th class="text-right">利息</th>
                <th class="text-right">剩余本金</th>
                <th>备注</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="detail in schedule" :key="detail.month" :class="{ 'extra-payment-row': detail.note }">
                <td class="period-cell">{{ detail.month }}</td>
                <td class="date-cell">{{ formatDate(detail.payment_date) }}</td>
                <td class="amount principal">¥{{ Math.round(detail.principal).toLocaleString('zh-CN') }}</td>
                <td class="amount interest">¥{{ Math.round(detail.interest).toLocaleString('zh-CN') }}</td>
                <td class="amount balance">¥{{ Math.round(detail.remaining_balance).toLocaleString('zh-CN') }}</td>
                <td class="note-cell">
                  <span v-if="detail.note" class="note-badge">{{ detail.note }}</span>
                </td>
              </tr>
            </tbody>
            <tfoot v-if="schedule.length > 0">
              <tr class="summary-row">
                <td colspan="2" class="summary-label">合计</td>
                <td class="amount summary-value">¥{{ Math.round(totalPrincipal).toLocaleString('zh-CN') }}</td>
                <td class="amount summary-value">¥{{ Math.round(totalInterest).toLocaleString('zh-CN') }}</td>
                <td colspan="2"></td>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-secondary" @click="handleClose">关闭</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { GetPaymentSchedule } from '../../wailsjs/go/main/App';

interface Props {
  show: boolean;
  loanId: number | null;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: 'close'): void;
}>();

interface PaymentDetail {
  month: number;
  payment_date: string;
  principal: number;
  interest: number;
  total_payment: number;
  remaining_balance: number;
  extra_payment: number;
  extra_interest: number;
  note: string;
}

const schedule = ref<PaymentDetail[]>([]);
const loading = ref(false);

const totalPrincipal = computed(() => {
  return schedule.value.reduce((sum, detail) => sum + detail.principal, 0);
});

const totalInterest = computed(() => {
  return schedule.value.reduce((sum, detail) => sum + detail.interest, 0);
});

watch(() => props.show, async (newVal) => {
  if (newVal && props.loanId) {
    loading.value = true;
    try {
      schedule.value = await GetPaymentSchedule(props.loanId);
    } catch (err) {
      console.error('获取还款明细失败:', err);
      alert('获取还款明细失败');
    } finally {
      loading.value = false;
    }
  }
});

const formatDate = (dateStr: string): string => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN');
};

const handleClose = () => {
  emit('close');
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: var(--glass-bg);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-lg);
  width: 90%;
  max-width: 1200px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid var(--glass-border);
}

.modal-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

.btn-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 0;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  transition: all var(--transition-base);
}

.btn-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: var(--text-secondary);
}

.schedule-table-wrapper {
  overflow-x: auto;
}

.schedule-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
  table-layout: fixed;
}

.schedule-table thead {
  background: rgba(255, 255, 255, 0.08);
  position: sticky;
  top: 0;
  z-index: 1;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.schedule-table th {
  padding: 0.625rem 0.75rem;
  font-weight: 600;
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  white-space: nowrap;
  border-bottom: 2px solid rgba(255, 255, 255, 0.15);
  text-align: center;
}

.schedule-table th:nth-child(1) { width: 8%; }   /* 期数 */
.schedule-table th:nth-child(2) { width: 14%; }  /* 还款日期 */
.schedule-table th:nth-child(3) { width: 20%; }  /* 本金 */
.schedule-table th:nth-child(4) { width: 20%; }  /* 利息 */
.schedule-table th:nth-child(5) { width: 22%; }  /* 剩余本金 */
.schedule-table th:nth-child(6) { width: 16%; }  /* 备注 */

.schedule-table th.text-right {
  text-align: right;
  padding-right: 1rem;
}

.schedule-table tbody tr {
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
  transition: background var(--transition-base);
}

.schedule-table tbody tr:hover {
  background: rgba(255, 255, 255, 0.04);
}

.schedule-table tbody tr.extra-payment-row {
  background: rgba(16, 185, 129, 0.08);
  border-left: 3px solid #10b981;
}

.schedule-table tbody tr.extra-payment-row:hover {
  background: rgba(16, 185, 129, 0.12);
}

.schedule-table td {
  padding: 0.5rem 0.75rem;
  color: var(--text-primary);
  white-space: nowrap;
}

.schedule-table td.period-cell {
  color: rgba(255, 255, 255, 0.6);
  font-weight: 600;
  text-align: center;
  font-size: 0.8125rem;
}

.schedule-table td.date-cell {
  color: rgba(255, 255, 255, 0.75);
  text-align: center;
  font-size: 0.8125rem;
}

.schedule-table td.amount {
  font-family: var(--font-mono);
  text-align: right;
  font-weight: 500;
  font-size: 0.875rem;
  letter-spacing: 0.01em;
  padding-right: 1rem;
}

.schedule-table td.principal {
  color: #3b82f6;
  font-weight: 600;
}

.schedule-table td.interest {
  color: #f59e0b;
  font-weight: 600;
}

.schedule-table td.balance {
  color: #10b981;
  font-weight: 700;
}

.schedule-table td.note-cell {
  text-align: center;
}

.note-badge {
  display: inline-block;
  padding: 0.25rem 0.625rem;
  font-size: 0.6875rem;
  font-weight: 600;
  color: #10b981;
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid rgba(16, 185, 129, 0.4);
  border-radius: 4px;
  letter-spacing: 0.025em;
}

.schedule-table tfoot {
  border-top: 2px solid rgba(255, 255, 255, 0.15);
  background: rgba(255, 255, 255, 0.05);
}

.schedule-table tfoot tr.summary-row {
  background: rgba(0, 128, 255, 0.08);
}

.schedule-table tfoot td {
  padding: 0.75rem 0.75rem;
}

.schedule-table tfoot td.summary-label {
  font-weight: 700;
  color: rgba(255, 255, 255, 0.95);
  text-align: right;
  padding-right: 1rem;
  font-size: 0.875rem;
  letter-spacing: 0.05em;
}

.schedule-table tfoot td.summary-value {
  font-weight: 700;
  font-size: 0.9375rem;
  font-family: var(--font-mono);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  padding: 1.5rem;
  border-top: 1px solid var(--glass-border);
}
</style>

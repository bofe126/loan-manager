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
                <th>本金</th>
                <th>利息</th>
                <th>剩余本金</th>
                <th>备注</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="detail in schedule" :key="detail.month">
                <td>{{ detail.month }}</td>
                <td>{{ formatDate(detail.payment_date) }}</td>
                <td class="amount principal">¥{{ detail.principal.toLocaleString('zh-CN', { minimumFractionDigits: 2 }) }}</td>
                <td class="amount interest">¥{{ detail.interest.toLocaleString('zh-CN', { minimumFractionDigits: 2 }) }}</td>
                <td class="amount">¥{{ detail.remaining_balance.toLocaleString('zh-CN', { minimumFractionDigits: 2 }) }}</td>
                <td>
                  <span v-if="detail.note" class="note-text">{{ detail.note }}</span>
                </td>
              </tr>
            </tbody>
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
import { ref, watch } from 'vue';
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
}

.schedule-table thead {
  background: rgba(255, 255, 255, 0.03);
}

.schedule-table th {
  padding: 0.75rem;
  font-weight: 600;
  font-size: 0.75rem;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  white-space: nowrap;
  border-bottom: 1px solid var(--glass-border);
  text-align: left;
}

.schedule-table tbody tr {
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  transition: all var(--transition-base);
}

.schedule-table tbody tr:hover {
  background: rgba(255, 255, 255, 0.03);
}

.schedule-table td {
  padding: 0.75rem;
  color: var(--text-primary);
  white-space: nowrap;
}

.schedule-table td.amount {
  font-family: var(--font-mono);
  text-align: right;
}

.schedule-table td.principal {
  color: #0080ff;
}

.schedule-table td.interest {
  color: #fbbf24;
}

.note-text {
  font-size: 0.75rem;
  color: #39ff14;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  padding: 1.5rem;
  border-top: 1px solid var(--glass-border);
}
</style>

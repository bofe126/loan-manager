<template>
  <div class="card">
    <div class="card-body">
      <div class="card-header-section">
        <h5 class="card-title">{{ title }}</h5>
        <div v-if="isEditMode && formData.remaining_amount > 0" class="remaining-amount-badge">
          贷款余额：¥{{ formData.remaining_amount.toLocaleString('zh-CN', { minimumFractionDigits: 2 }) }}
        </div>
      </div>
      <form @submit.prevent="handleSubmit">
        <div class="form-grid">
          <div class="form-group">
            <label class="form-label">借款人</label>
            <input
              v-model="formData.borrower_name"
              type="text"
              class="form-control"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label">银行名称</label>
            <input
              v-model="formData.bank_name"
              type="text"
              class="form-control"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label">贷款总金额（元）</label>
            <input
              v-model.number="formData.total_amount"
              type="number"
              step="0.01"
              class="form-control"
              required
              @input="updatePreview"
            />
          </div>

          <div class="form-group">
            <label class="form-label">年利率（%）</label>
            <input
              v-model.number="formData.interest_rate"
              type="number"
              step="0.01"
              class="form-control"
              required
              @input="updatePreview"
            />
          </div>

          <div class="form-group">
            <label class="form-label">开始日期</label>
            <input
              v-model="formData.start_date"
              type="date"
              class="form-control"
              required
              @change="updatePreview"
            />
          </div>

          <div class="form-group">
            <label class="form-label">结束日期</label>
            <input
              v-model="formData.end_date"
              type="date"
              class="form-control"
              required
              @change="updatePreview"
            />
          </div>

          <div class="form-group">
            <label class="form-label">第一期还款日期</label>
            <input
              v-model="formData.first_payment_date"
              type="date"
              class="form-control"
              placeholder="可选，默认为开始日期的下个月"
            />
          </div>

          <div class="form-group">
            <label class="form-label">每月还款日</label>
            <input
              v-model.number="formData.payment_date"
              type="number"
              min="1"
              max="31"
              class="form-control"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label">还款方式</label>
            <select
              v-model="formData.payment_method"
              class="form-select"
              required
              @change="updatePreview"
            >
              <option value="等额本息">等额本息</option>
              <option value="等额本金">等额本金</option>
              <option value="先息后本">先息后本</option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label">状态</label>
            <select v-model="formData.status" class="form-select" required>
              <option value="active">进行中</option>
              <option value="completed">已完成</option>
              <option value="overdue">逾期</option>
            </select>
          </div>
        </div>

        <div v-if="previewMonthlyPayment > 0" class="preview-box">
          <span class="preview-label">{{ isEditMode ? '当前月供' : '预计月供' }}</span>
          <span class="preview-value">¥{{ Math.round(previewMonthlyPayment).toLocaleString('zh-CN') }}</span>
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-success" @click="handleSubmitAndNew" :disabled="loading">
            {{ loading ? '提交中...' : '再记一条' }}
          </button>
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '提交中...' : '提交' }}
          </button>
          <button type="button" class="btn btn-secondary" @click="handleCancel">
            取消
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { models } from '../../wailsjs/go/models';
import { useLoanCalculator } from '../composables/useLoanCalculator';

interface Props {
  title: string;
  initialData?: models.Loan;
  loading?: boolean;
}

interface Emits {
  (e: 'submit', data: Partial<models.Loan>): void;
  (e: 'submitAndNew', data: Partial<models.Loan>): void;
  (e: 'cancel'): void;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
});

const emit = defineEmits<Emits>();

const { calculateMonthlyPayment } = useLoanCalculator();

const formData = ref({
  borrower_name: '',
  bank_name: '',
  total_amount: 0,
  remaining_amount: 0,
  interest_rate: 0,
  start_date: '',
  end_date: '',
  first_payment_date: '',
  payment_date: 1,
  payment_method: '等额本息',
  status: 'active',
  monthly_payment: 0,
});

const previewMonthlyPayment = ref(0);
const isEditMode = ref(false);

// 初始化表单数据
onMounted(() => {
  if (props.initialData) {
    isEditMode.value = true;
    const startDate = props.initialData.start_date ? new Date(props.initialData.start_date).toISOString().split('T')[0] : '';
    const endDate = props.initialData.end_date ? new Date(props.initialData.end_date).toISOString().split('T')[0] : '';
    const firstPaymentDate = props.initialData.first_payment_date ? new Date(props.initialData.first_payment_date).toISOString().split('T')[0] : '';

    formData.value = {
      borrower_name: props.initialData.borrower_name,
      bank_name: props.initialData.bank_name,
      total_amount: props.initialData.total_amount || 0,
      remaining_amount: props.initialData.remaining_amount || 0,
      interest_rate: props.initialData.interest_rate,
      start_date: startDate,
      end_date: endDate,
      first_payment_date: firstPaymentDate,
      payment_date: props.initialData.payment_date,
      payment_method: props.initialData.payment_method,
      status: props.initialData.status,
      monthly_payment: props.initialData.monthly_payment || 0,
    };
    // 编辑模式下直接使用后端返回的月供
    previewMonthlyPayment.value = props.initialData.monthly_payment || 0;
    console.log('=== 前端初始化 ===');
    console.log('后端返回的月供:', props.initialData.monthly_payment);
    console.log('设置 previewMonthlyPayment 为:', previewMonthlyPayment.value);
    console.log('isEditMode:', isEditMode.value);
  }
});

// 更新预览
const updatePreview = () => {
  console.log('updatePreview 被调用, isEditMode:', isEditMode.value);

  // 编辑模式下不重新计算，使用后端返回的值
  if (isEditMode.value) {
    console.log('编辑模式，不重新计算，保持月供:', previewMonthlyPayment.value);
    return;
  }

  // 新建模式下才计算预览
  if (
    formData.value.total_amount > 0 &&
    formData.value.interest_rate > 0 &&
    formData.value.start_date &&
    formData.value.end_date
  ) {
    const startDate = new Date(formData.value.start_date);
    const endDate = new Date(formData.value.end_date);

    const calculated = calculateMonthlyPayment(
      formData.value.total_amount,
      formData.value.interest_rate,
      startDate,
      endDate,
      formData.value.payment_method as any
    );

    console.log('新建模式，计算月供:', calculated);
    previewMonthlyPayment.value = calculated;
  }
};

// 提交表单
const handleSubmit = () => {
  console.log('handleSubmit called');
  const submitData: any = {
    ...formData.value,
    start_date: new Date(formData.value.start_date).toISOString(),
    end_date: new Date(formData.value.end_date).toISOString(),
  };

  // 如果有第一期还款日期，转换为ISO格式；否则删除该字段
  if (formData.value.first_payment_date && formData.value.first_payment_date.trim() !== '') {
    submitData.first_payment_date = new Date(formData.value.first_payment_date).toISOString();
  } else {
    delete submitData.first_payment_date;
  }

  if (props.initialData) {
    submitData.id = props.initialData.id;
  }

  console.log('Emitting submit event with data:', submitData);
  emit('submit', submitData);
};

// 提交并新建
const handleSubmitAndNew = () => {
  console.log('handleSubmitAndNew called');
  const submitData: any = {
    ...formData.value,
    start_date: new Date(formData.value.start_date).toISOString(),
    end_date: new Date(formData.value.end_date).toISOString(),
  };

  // 如果有第一期还款日期，转换为ISO格式；否则删除该字段
  if (formData.value.first_payment_date && formData.value.first_payment_date.trim() !== '') {
    submitData.first_payment_date = new Date(formData.value.first_payment_date).toISOString();
  } else {
    delete submitData.first_payment_date;
  }

  if (props.initialData) {
    submitData.id = props.initialData.id;
  }

  console.log('Emitting submitAndNew event with data:', submitData);
  emit('submitAndNew', submitData);
};

// 取消
const handleCancel = () => {
  emit('cancel');
};
</script>

<style scoped>
.card-header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.card-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

.remaining-amount-badge {
  padding: 0.5rem 1rem;
  background: rgba(0, 128, 255, 0.1);
  border: 1px solid rgba(0, 128, 255, 0.3);
  border-radius: var(--radius-md);
  font-size: 0.9375rem;
  font-weight: 600;
  color: #0080ff;
  font-family: var(--font-mono);
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-label {
  font-size: 0.8125rem;
  margin-bottom: 0.375rem;
}

.form-control,
.form-select {
  padding: 0.5rem 0.75rem;
  font-size: 0.9375rem;
}

.preview-box {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.25rem;
  background: rgba(0, 128, 255, 0.1);
  border: 1px solid rgba(0, 128, 255, 0.3);
  border-radius: var(--radius-md);
  margin-bottom: 1.5rem;
}

.preview-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.preview-value {
  font-size: 1.25rem;
  font-weight: 700;
  color: #0080ff;
  font-family: var(--font-mono);
}

.form-actions {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}

.form-actions .btn {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
}

@media (max-width: 1024px) {
  .form-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>

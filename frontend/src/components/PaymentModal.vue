<template>
  <div
    class="modal fade"
    :class="{ show: show }"
    :style="{ display: show ? 'block' : 'none' }"
    tabindex="-1"
  >
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">记录还款</h5>
          <button
            type="button"
            class="btn-close"
            @click="$emit('close')"
          ></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="handleSubmit">
            <div class="mb-3">
              <label class="form-label">还款金额（元）</label>
              <input
                v-model.number="amount"
                type="number"
                step="0.01"
                class="form-control"
                required
                min="0.01"
              />
            </div>
            <div v-if="error" class="alert alert-danger">{{ error }}</div>
          </form>
        </div>
        <div class="modal-footer">
          <button
            type="button"
            class="btn btn-secondary"
            @click="$emit('close')"
          >
            取消
          </button>
          <button
            type="button"
            class="btn btn-primary"
            @click="handleSubmit"
            :disabled="loading"
          >
            {{ loading ? '提交中...' : '确认还款' }}
          </button>
        </div>
      </div>
    </div>
  </div>
  <div
    v-if="show"
    class="modal-backdrop fade show"
    @click="$emit('close')"
  ></div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

interface Props {
  show: boolean;
  loanId: number | null;
  loading?: boolean;
}

interface Emits {
  (e: 'close'): void;
  (e: 'submit', loanId: number, amount: number): void;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
});

const emit = defineEmits<Emits>();

const amount = ref(0);
const error = ref<string | null>(null);

watch(
  () => props.show,
  (newVal) => {
    if (newVal) {
      amount.value = 0;
      error.value = null;
    }
  }
);

const handleSubmit = () => {
  error.value = null;

  if (!props.loanId) {
    error.value = '无效的贷款ID';
    return;
  }

  if (amount.value <= 0) {
    error.value = '还款金额必须大于0';
    return;
  }

  emit('submit', props.loanId, amount.value);
};
</script>

<style scoped>
.modal.show {
  display: block;
}
</style>

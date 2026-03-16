<template>
  <div class="container mt-4">
    <LoanList
      :loans="loans"
      :loading="loading"
      @edit="handleEdit"
      @delete="handleDelete"
      @payment="handlePayment"
    />

    <PaymentModal
      :show="showPaymentModal"
      :loan-id="selectedLoanId"
      :loading="paymentLoading"
      @close="showPaymentModal = false"
      @submit="handlePaymentSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useLoans } from '../composables/useLoans';
import { MakePayment } from '../../wailsjs/go/main/App';
import LoanList from '../components/LoanList.vue';
import PaymentModal from '../components/PaymentModal.vue';

const router = useRouter();
const { loans, loading, loadLoans, deleteLoan } = useLoans();

const showPaymentModal = ref(false);
const selectedLoanId = ref<number | null>(null);
const paymentLoading = ref(false);

onMounted(async () => {
  await loadLoans();
});

const handleEdit = (id: number) => {
  router.push(`/edit/${id}`);
};

const handleDelete = async (id: number) => {
  await deleteLoan(id);
};

const handlePayment = (id: number) => {
  selectedLoanId.value = id;
  showPaymentModal.value = true;
};

const handlePaymentSubmit = async (loanId: number, amount: number, paymentDate: string) => {
  paymentLoading.value = true;
  try {
    // 将日期转换为 ISO 8601 格式
    const dateObj = new Date(paymentDate);
    const isoDate = dateObj.toISOString();
    await MakePayment(loanId, amount, isoDate);
    showPaymentModal.value = false;
    await loadLoans();
  } catch (err) {
    console.error('还款失败:', err);
    alert('还款失败，请重试');
  } finally {
    paymentLoading.value = false;
  }
};
</script>

<style scoped>
/* 样式已移除，页面更简洁 */
</style>

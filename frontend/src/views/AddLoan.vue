<template>
  <div class="container mt-4">
    <LoanForm
      title="新建贷款"
      :loading="loading"
      @submit="handleSubmit"
      @submitAndNew="handleSubmitAndNew"
      @cancel="handleCancel"
    />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useLoans } from '../composables/useLoans';
import LoanForm from '../components/LoanForm.vue';
import { models } from '../../wailsjs/go/models';

const router = useRouter();
const { loading, createLoan } = useLoans();

const handleSubmit = async (data: Partial<models.Loan>) => {
  console.log('AddLoan handleSubmit called with data:', data);
  const success = await createLoan(data);
  console.log('createLoan result:', success);
  if (success) {
    router.push('/');
  }
};

const handleSubmitAndNew = async (data: Partial<models.Loan>) => {
  console.log('AddLoan handleSubmitAndNew called with data:', data);
  const success = await createLoan(data);
  console.log('createLoan result:', success);
  if (success) {
    // 刷新页面以清空表单
    window.location.reload();
  }
};

const handleCancel = () => {
  router.push('/');
};
</script>

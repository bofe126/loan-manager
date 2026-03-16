<template>
  <div class="container mt-4">
    <h2 class="mb-4">编辑贷款</h2>
    <div v-if="loading" class="text-center">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">加载中...</span>
      </div>
    </div>
    <LoanForm
      v-else-if="loan"
      title="编辑贷款"
      :initial-data="loan"
      :loading="updating"
      @submit="handleSubmit"
      @cancel="handleCancel"
    />
    <div v-else class="alert alert-danger">贷款不存在</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useLoans } from '../composables/useLoans';
import LoanForm from '../components/LoanForm.vue';
import { models } from '../../wailsjs/go/models';

const router = useRouter();
const route = useRoute();
const { loading, getLoan, updateLoan } = useLoans();

const loan = ref<models.Loan | null>(null);
const updating = ref(false);

onMounted(async () => {
  const id = Number(route.params.id);
  if (id) {
    loan.value = await getLoan(id);
  }
});

const handleSubmit = async (data: Partial<models.Loan>) => {
  updating.value = true;
  const success = await updateLoan(data as models.Loan);
  updating.value = false;
  if (success) {
    router.push('/');
  }
};

const handleCancel = () => {
  router.push('/');
};
</script>

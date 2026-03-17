<template>
  <div class="row mb-4">
    <div class="col-md-6">
      <div class="card text-white bg-primary stats-card">
        <div class="card-body">
          <div class="stats-header">
            <h5 class="card-title">月还款</h5>
          </div>
          <h2 class="stats-value">
            <span class="currency">¥</span>
            {{ Math.round(totalMonthlyPayment).toLocaleString('zh-CN') }}
          </h2>
          <div class="stats-footer">
            <span class="stats-label">MONTHLY PAYMENT</span>
          </div>
        </div>
      </div>
    </div>
    <div class="col-md-6">
      <div class="card text-white bg-success stats-card">
        <div class="card-body">
          <div class="stats-header">
            <h5 class="card-title">剩余总额</h5>
          </div>
          <h2 class="stats-value">
            <span class="currency">¥</span>
            {{ Math.round(totalRemainingAmount).toLocaleString('zh-CN') }}
          </h2>
          <div class="stats-footer">
            <span class="stats-label">REMAINING BALANCE</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { models } from '../../wailsjs/go/models';
import { useLoanCalculator } from '../composables/useLoanCalculator';

interface Props {
  loans: models.Loan[];
}

const props = defineProps<Props>();

const { calculateRemainingAmount } = useLoanCalculator();

const totalMonthlyPayment = computed(() => {
  return props.loans.reduce((sum, loan) => {
    if (loan.status === 'active') {
      return sum + loan.monthly_payment;
    }
    return sum;
  }, 0);
});

const totalRemainingAmount = computed(() => {
  return props.loans.reduce((sum, loan) => {
    if (loan.status === 'active') {
      return sum + loan.remaining_amount;
    }
    return sum;
  }, 0);
});
</script>

<style scoped>
.stats-card {
  position: relative;
  overflow: hidden;
  min-height: 180px;
}

.stats-card .card-body {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
}

.stats-header {
  margin-bottom: 1rem;
}

.stats-header .card-title {
  margin: 0;
  font-size: 0.875rem;
  font-weight: 400;
  opacity: 0.9;
  text-transform: uppercase;
  letter-spacing: 0.1em;
}

.stats-value {
  font-family: var(--font-heading);
  font-size: 2.5rem;
  font-weight: 300;
  margin: 1rem 0;
  text-shadow: 0 2px 12px rgba(0, 0, 0, 0.3);
  letter-spacing: -0.02em;
  line-height: 1;
}

.currency {
  font-size: 1.5rem;
  opacity: 0.8;
  margin-right: 0.25rem;
}

.stats-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
}

.stats-label {
  font-family: var(--font-heading);
  font-size: 0.625rem;
  font-weight: 400;
  letter-spacing: 0.15em;
  opacity: 0.7;
}
</style>

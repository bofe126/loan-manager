<template>
  <div class="home">
    <div class="dashboard-layout">
      <!-- 左侧统计卡片 -->
      <div class="stats-sidebar">
        <div class="stat-card" @click="handleMonthlyPaymentClick">
          <div class="stat-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-label">月还款总额</div>
            <div class="stat-value">¥{{ Math.round(totalMonthlyPayment).toLocaleString('zh-CN') }}</div>
          </div>
        </div>

        <div class="stat-card" @click="handleRemainingAmountClick">
          <div class="stat-icon primary">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M3 3v18h18"/>
              <path d="m19 9-5 5-4-4-3 3"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-label">剩余贷款</div>
            <div class="stat-value">¥{{ Math.round(totalRemainingAmount).toLocaleString('zh-CN') }}</div>
          </div>
        </div>

        <div class="stat-card" @click="handleLoanCountClick">
          <div class="stat-icon warning">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect width="18" height="18" x="3" y="3" rx="2"/>
              <path d="M3 9h18"/>
              <path d="M9 21V9"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-label">贷款数量</div>
            <div class="stat-value">{{ activeLoans }}</div>
          </div>
        </div>
      </div>

      <!-- 右侧图表区域 -->
      <div class="charts-wrapper">
        <TrendCharts ref="trendChartsRef" :loans="loans" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useLoans } from '../composables/useLoans';
import { useLoanCalculator } from '../composables/useLoanCalculator';
import TrendCharts from '../components/TrendCharts.vue';

const router = useRouter();
const { loans, loadLoans } = useLoans();
const { calculateRemainingAmount } = useLoanCalculator();

const trendChartsRef = ref<InstanceType<typeof TrendCharts> | null>(null);

let refreshInterval: number | null = null;

onMounted(async () => {
  await loadLoans();
  refreshInterval = window.setInterval(() => {
    loadLoans();
  }, 60000);
});

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval);
  }
});

const totalMonthlyPayment = computed(() => {
  return loans.value.reduce((sum, loan) => {
    if (loan.status === 'active') {
      return sum + loan.monthly_payment;
    }
    return sum;
  }, 0);
});

const totalRemainingAmount = computed(() => {
  return loans.value.reduce((sum, loan) => {
    if (loan.status === 'active') {
      return sum + loan.remaining_amount;
    }
    return sum;
  }, 0);
});

const activeLoans = computed(() => loans.value.filter(l => l.status === 'active').length);

// 点击处理函数
const handleMonthlyPaymentClick = () => {
  if (trendChartsRef.value) {
    (trendChartsRef.value as any).chartType = 'monthly';
  }
};

const handleRemainingAmountClick = () => {
  if (trendChartsRef.value) {
    (trendChartsRef.value as any).chartType = 'remaining';
  }
};

const handleLoanCountClick = () => {
  router.push('/loans');
};
</script>

<style scoped>
.home {
  min-height: calc(100vh - 56px);
  padding: 1.5rem;
}

/* 仪表板布局 */
.dashboard-layout {
  display: flex;
  gap: 1.5rem;
  align-items: stretch;
}

/* 左侧统计卡片 - 黄金分割比例 */
.stats-sidebar {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  flex: 0 0 38.2%;
  max-width: 38.2%;
}

.stat-card {
  background: var(--glass-bg);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-md);
  padding: 1rem;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 1rem;
  transition: all var(--transition-base);
  position: relative;
  overflow: hidden;
  flex: 1;
  min-height: 80px;
  cursor: pointer;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.05), transparent);
  transition: left 0.5s;
}

.stat-card:hover::before {
  left: 100%;
}

.stat-card:hover {
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.4);
}

.stat-card:active {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.3);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-sm);
  background: linear-gradient(135deg, rgba(0, 128, 255, 0.2) 0%, rgba(51, 153, 255, 0.2) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #0080ff;
  flex-shrink: 0;
}

.stat-icon svg {
  width: 24px;
  height: 24px;
}

.stat-icon.primary {
  background: linear-gradient(135deg, rgba(0, 128, 255, 0.2) 0%, rgba(51, 153, 255, 0.2) 100%);
  color: #0080ff;
}

.stat-icon.success {
  background: linear-gradient(135deg, rgba(57, 255, 20, 0.2) 0%, rgba(46, 204, 64, 0.2) 100%);
  color: #39ff14;
}

.stat-icon.warning {
  background: linear-gradient(135deg, rgba(247, 147, 26, 0.2) 0%, rgba(255, 165, 0, 0.2) 100%);
  color: #F7931A;
}

.stat-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.25rem;
}

.stat-label {
  font-size: 0.75rem;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 500;
}

.stat-value {
  font-size: 1.375rem;
  font-weight: 700;
  color: var(--text-primary);
  font-family: var(--font-mono);
  line-height: 1.2;
}

/* 右侧图表区域 - 黄金分割比例 */
.charts-wrapper {
  flex: 0 0 61.8%;
  max-width: 61.8%;
  min-width: 0;
  display: flex;
  flex-direction: column;
  padding-right: 1.5rem;
}

/* 响应式 */
@media (max-width: 768px) {
  .home {
    padding: 1rem;
  }

  .dashboard-layout {
    flex-direction: column;
  }

  .stats-sidebar {
    width: 100%;
    flex-direction: row;
    overflow-x: auto;
  }

  .stat-card {
    min-width: 160px;
  }
}
</style>

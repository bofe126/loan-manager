<template>
  <div class="chart-container">
    <div class="card chart-card">
      <div class="card-body">
        <div class="chart-header">
          <h5 class="card-title">{{ currentChartTitle }}</h5>
          <div class="chart-controls">
            <button
              class="chart-toggle primary-btn"
              :class="{ active: chartType === 'monthly' }"
              @click="chartType = 'monthly'"
              data-type="monthly"
            >
              月还款
            </button>
            <button
              class="chart-toggle secondary-btn"
              :class="{ active: chartType === 'remaining' }"
              @click="chartType = 'remaining'"
              data-type="remaining"
            >
              剩余贷款
            </button>
          </div>
        </div>
        <canvas ref="chartCanvas"></canvas>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { Chart, registerables } from 'chart.js';
import { models } from '../../wailsjs/go/models';
import { useTrendData } from '../composables/useTrendData';

Chart.register(...registerables);

interface Props {
  loans: models.Loan[];
}

const props = defineProps<Props>();

const { calculateTrendData } = useTrendData();

const chartCanvas = ref<HTMLCanvasElement | null>(null);
const chartType = ref<'monthly' | 'remaining'>('monthly');

// 暴露 chartType 给父组件
defineExpose({
  chartType
});

let currentChart: Chart | null = null;

const currentChartTitle = computed(() => {
  return chartType.value === 'monthly' ? '月还款趋势' : '剩余贷款趋势';
});

// 金融风格图表配置
const createFinancialChartConfig = (
  labels: string[],
  data: number[],
  label: string,
  color: { border: string; gradient: string[] }
) => {
  return {
    type: 'line' as const,
    data: {
      labels,
      datasets: [
        {
          label,
          data,
          borderColor: color.border,
          borderWidth: 3,
          backgroundColor: (context: any) => {
            const ctx = context.chart.ctx;
            const gradient = ctx.createLinearGradient(0, 0, 0, 350);
            gradient.addColorStop(0, color.gradient[0]);
            gradient.addColorStop(0.5, color.gradient[1]);
            gradient.addColorStop(1, color.gradient[2]);
            return gradient;
          },
          fill: true,
          tension: 0.4,
          pointRadius: 0,
          pointHoverRadius: 8,
          pointHoverBackgroundColor: color.border,
          pointHoverBorderColor: '#fff',
          pointHoverBorderWidth: 3,
          pointHitRadius: 30,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: true,
      interaction: {
        intersect: false,
        mode: 'index' as const,
      },
      plugins: {
        legend: {
          display: false,
        },
        tooltip: {
          enabled: true,
          backgroundColor: 'rgba(10, 14, 39, 0.95)',
          titleColor: '#FFFFFF',
          bodyColor: '#7B96FF',
          borderColor: 'rgba(98, 126, 234, 0.3)',
          borderWidth: 1,
          padding: 16,
          displayColors: false,
          titleFont: {
            size: 12,
            weight: 'bold' as const,
            family: "'Inter', sans-serif",
          },
          bodyFont: {
            size: 18,
            weight: 'bold' as const,
            family: "'JetBrains Mono', monospace",
          },
          callbacks: {
            label: (context: any) => {
              return `¥ ${Math.round(context.parsed.y).toLocaleString('zh-CN')}`;
            },
          },
          cornerRadius: 8,
          caretSize: 8,
        },
      },
      scales: {
        x: {
          display: false,
        },
        y: {
          display: false,
        },
      },
    },
  };
};

const renderChart = () => {
  if (!chartCanvas.value) return;

  const trendData = calculateTrendData(props.loans as any);

  // 销毁旧图表
  if (currentChart) {
    currentChart.destroy();
  }

  // 根据当前选择的图表类型渲染
  if (chartType.value === 'monthly') {
    const monthlyConfig = createFinancialChartConfig(
      trendData.labels,
      trendData.monthlyPayments,
      '月还款',
      {
        border: '#39ff14',
        gradient: [
          'rgba(57, 255, 20, 0.5)',
          'rgba(57, 255, 20, 0.2)',
          'rgba(57, 255, 20, 0.01)'
        ],
      }
    );
    currentChart = new Chart(chartCanvas.value, monthlyConfig);
  } else {
    const remainingConfig = createFinancialChartConfig(
      trendData.labels,
      trendData.remainingAmounts,
      '剩余贷款',
      {
        border: '#0080ff',
        gradient: [
          'rgba(0, 128, 255, 0.5)',
          'rgba(0, 128, 255, 0.2)',
          'rgba(0, 128, 255, 0.01)'
        ],
      }
    );
    currentChart = new Chart(chartCanvas.value, remainingConfig);
  }
};

onMounted(() => {
  renderChart();
});

watch(
  () => props.loans,
  () => {
    renderChart();
  },
  { deep: true }
);

watch(chartType, () => {
  renderChart();
});
</script>

<style scoped>
.chart-container {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
}

.chart-card {
  position: relative;
  overflow: hidden;
  background: var(--glass-bg);
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius-lg);
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.4), inset 0 1px 0 0 rgba(255, 255, 255, 0.1);
}

.chart-card::after {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(98, 126, 234, 0.08) 0%, transparent 70%);
  pointer-events: none;
  animation: pulse 4s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 0.5;
    transform: scale(1) rotate(0deg);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.1) rotate(180deg);
  }
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 0;
  border-bottom: none;
  position: relative;
  z-index: 2;
}

.chart-header .card-title {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: 0.025em;
  font-family: var(--font-heading);
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.chart-controls {
  display: flex;
  gap: 0.5rem;
  background: rgba(17, 24, 39, 0.5);
  padding: 0.25rem;
  border-radius: var(--radius-md);
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.chart-toggle {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 0.8125rem;
  font-weight: 500;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  cursor: pointer;
  padding: 0.375rem 1rem;
  border-radius: calc(var(--radius-md) - 0.25rem);
  transition: all var(--transition-base);
  font-family: var(--font-body);
  position: relative;
  overflow: hidden;
  display: inline-flex;
  align-items: center;
}

.chart-toggle.primary-btn {
  background: transparent;
  color: var(--text-secondary);
  font-weight: 500;
  padding: 0.375rem 1rem;
}

.chart-toggle.primary-btn:hover {
  color: #39ff14;
}

.chart-toggle.primary-btn.active {
  background: linear-gradient(135deg, #39ff14 0%, #2ecc40 100%);
  color: #0a0e27;
  font-weight: 600;
  box-shadow: 0 0 20px rgba(57, 255, 20, 0.6), inset 0 1px 0 0 rgba(255, 255, 255, 0.2);
}

.chart-toggle.secondary-btn {
  background: transparent;
  color: var(--text-secondary);
  font-weight: 500;
  padding: 0.375rem 1rem;
}

.chart-toggle.secondary-btn:hover {
  color: #0080ff;
}

.chart-toggle.secondary-btn.active {
  background: linear-gradient(135deg, #0080ff 0%, #0066cc 100%);
  color: white;
  font-weight: 600;
  box-shadow: 0 0 20px rgba(0, 128, 255, 0.6), inset 0 1px 0 0 rgba(255, 255, 255, 0.2);
}

.chart-toggle::before {
  display: none;
}

.chart-toggle:hover {
  color: var(--text-primary);
}

canvas {
  max-height: 350px;
  position: relative;
  z-index: 1;
}

@media (max-width: 768px) {
  .chart-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .chart-controls {
    width: 100%;
  }

  .chart-toggle {
    flex: 1;
    padding: 0.625rem 1rem;
  }

  canvas {
    max-height: 280px;
  }
}
</style>

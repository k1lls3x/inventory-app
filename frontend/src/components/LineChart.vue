<template>
  <div style="height: 280px; padding: 1rem;">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>

<script setup>
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale
} from 'chart.js'
import { computed } from 'vue'

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale
)

const props = defineProps({
  data: {
    type: Object,
    required: true
  }
})

const chartData = computed(() => props.data)

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      callbacks: {
        label: ctx => `${ctx.formattedValue} шт.`
      }
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: {
        callback: v => `${v} шт.`,
        color: '#6b7280'
      },
      grid: { color: '#e5e7eb' }
    },
    x: {
      ticks: { color: '#6b7280', maxRotation: 0 },
      grid: { display: false }
    }
  }
}
</script>

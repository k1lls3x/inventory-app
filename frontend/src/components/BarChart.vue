<template>
  <div style="height: 280px; padding: 1rem;">
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>

<script setup>
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale
} from 'chart.js'
import { computed } from 'vue'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

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
      ticks: { maxRotation: 0, minRotation: 0 },
      grid: { display: false }
    }
  }
}
</script>

<template>
  <div style="height: 280px; padding: 1rem;">
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>

<script setup>
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale
} from 'chart.js'
import { Bar } from 'vue-chartjs'
import { computed } from 'vue'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

const props = defineProps({
  data: Array
})

const chartData = computed(() => ({
  labels: props.data.map(item => item.name),
  datasets: [
  {
    label: 'Поступило',
    data: props.data.map(item => item.received),
    backgroundColor: '#000',
    barThickness: 20
  }
]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    },
    title: {
      display: true,
      align: 'start',
      color: '#111827',
      font: {
        size: 16,
        weight: 'bold'
      },
      padding: {
        bottom: 20
      }
    },
    tooltip: {
      callbacks: {
        label: function (context) {
         return `${context.formattedValue.toLocaleString()} шт.`
        }
      }
    }
  },
  layout: {
    padding: {
      top: 10,
      bottom: 10,
      left: 10,
      right: 10
    }
  },
  scales: {
    x: {
  ticks: {
    maxRotation: 0,
    minRotation: 0
  }
},
    y: {
      beginAtZero: true,
      grid: {
        color: '#e5e7eb'
      },
      ticks: {
        color: '#6b7280',
        font: {
          size: 12
        },
         callback: (value) => `${value.toLocaleString()} шт.`
      }
    }
  }
}
</script>

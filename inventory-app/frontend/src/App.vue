<template>
  <div class="dashboard">
    <header class="header">
      <h1>📊 Складская система</h1>
      <nav class="tabs">
        <button
          v-for="tab in tabs"
          :key="tab"
          :class="{ active: currentTab === tab }"
          @click="currentTab = tab"
        >
          {{ tab }}
        </button>
      </nav>
    </header>

    <main>
      <section v-if="currentTab === 'Дашборд'">
        <section class="cards">
          <div class="card highlight">
            <p class="title">Всего остатков</p>
            <p class="value">{{ totalStock }}</p>
            <p class="note positive">+15% за месяц</p>
          </div>
          <div class="card">
            <p class="title">Товаров</p>
            <p class="value">{{ itemCount }}</p>
            <p class="note" v-if="newItems > 0">+{{ newItems }} новых за месяц</p>
           </div>
          <div class="card">
            <p class="title">Поставки</p>
            <p class="value">{{ monthlyOrders }}</p>
            <p class="note">в этом месяце</p>
          </div>
        </section>

        <section class="dashboard-visuals">
        <div class="chart-card">
          <p class="title">Остатки за неделю</p>
          <img src="https://fakeimg.pl/600x200/ddd/000/?text=LineChart" alt="Line chart" />
        </div>
        <div class="chart-card">
          <p class="title">Оборот по складам</p>
          <BarChart :data="turnoverData" />
        </div>
  </section>

        <section class="table-section">
          <p class="title">Популярные товары</p>
          <table>
            <thead>
              <tr>
                <th>Наименование</th>
                <th>SKU</th>
                <th>Склад</th>
                <th>Остаток</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="item in topItems"
                :key="item.id"
                :class="{ 'zero-stock': item.quantity === 0 }"
              >
                <td>{{ item.name }}</td>
                <td>{{ item.sku }}</td>
                <td>{{ item.warehouse }}</td>
                <td>{{ item.quantity }}</td>
              </tr>
            </tbody>
          </table>
        </section>
      </section>

      <section v-else>
        <p>Раздел "{{ currentTab }}" в разработке...</p>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { GetDashboard, GetTopItems, GetTurnoverByWarehouse } from '../wailsjs/go/app/App'
import BarChart from './components/BarChart.vue'

const tabs = [
  'Дашборд',
  'Остатки',
  'Поставки',
  'Товары',
  'Склады',
  'Поставщики',
  'Движения'
]

const currentTab = ref('Дашборд')

// Динамические данные
const totalStock = ref(0)
const itemCount = ref(0)
const monthlyOrders = ref(0)
const topItems = ref([])

const turnoverData = ref([]) // ⬅️ Добавлено: данные по обороту

const newItems = ref(0)

onMounted(() => {
  GetDashboard().then(data => {
    totalStock.value = data.total_stock
    itemCount.value = data.item_count
    monthlyOrders.value = data.monthly_orders
    newItems.value = data.new_items
  }).catch(err => {
    console.error("Ошибка загрузки дашборда:", err)
  })

  GetTopItems().then(data => {
    topItems.value = data
  }).catch(err => {
    console.error("Ошибка загрузки топовых товаров:", err)
  })

  GetTurnoverByWarehouse().then(data => {
    turnoverData.value = data
    console.log("Оборот по складам:", data)
  }).catch(err => {
    console.error("Ошибка загрузки оборота по складу:", err)
  })
})
</script>

<style scoped>
.dashboard {
  padding: 2rem;
  background-color: #f9f9f9;
}
.dashboard-visuals {
  display: flex;
  flex-wrap: wrap;
  gap: 1.5rem;
  margin-bottom: 2rem;
  align-items: stretch;
}

.header {
  margin-bottom: 2rem;
}

h1 {
  font-size: 2rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 1rem;
}

.tabs button {
  padding: 0.6rem 1.2rem;
  font-size: 0.95rem;
  font-weight: 700;
}

.card .title {
  font-size: 1rem;
  font-weight: 600;
  color: #4b5563;
}

.card .value {
  font-size: 2.2rem;
  font-weight: 800;
  color: #111827;
}

.card .note {
  font-size: 0.9rem;
  font-weight: 500;
  color: #6b7280;
}

.table-section .title {
  font-size: 1.2rem;
  font-weight: 700;
  margin-bottom: 1rem;
}

.table-section th {
  font-weight: 700;
  font-size: 0.95rem;
  color: #374151;
}

.table-section td {
  font-size: 0.9rem;
  font-weight: 500;
}

.tabs button:hover {
  background-color: #f3f4f6;
  transform: scale(1.04);
}
.tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.tabs button {
  padding: 0.5rem 1rem;
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.95rem;
  transition: background-color 0.15s ease, transform 0.15s ease;
  box-shadow: none;
  backface-visibility: hidden;
  transform: translateZ(0);
}

.tabs button.active {
  background-color: #2563eb;
  color: #fff;
  border-color: #2563eb;
}

.cards {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.card {
  flex: 1;
  background: white;
  border-radius: 10px;
  padding: 1.5rem;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
}

.card.highlight {
  border-left: 5px solid #2563eb;
}

.card .title {
  font-weight: 700;
  color: #1f2937; /* почти чёрный */
}

.card .value {
  font-size: 1.9rem;
  font-weight: bold;
  margin: 0.2rem 0;
}

.card .note {
  font-size: 0.8rem;
  color: #888;
}

.card .note.positive {
  color: #22c55e;
}

.charts {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.chart-card {
  flex: 1 1 45%;
  min-width: 320px;
  max-width: 100%;
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid #e5e7eb;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transition: all 0.2s ease-in-out;
}


.chart-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
}
.chart-card img {
  width: 100%;
  height: 200px;
  border-radius: 6px;
  object-fit: cover;
}

.table-section {
  background: white;
  border-radius: 10px;
  padding: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.table-section table {
  width: 100%;
  border-collapse: collapse;
}

.table-section th,
.table-section td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.table-section th {
  background-color: #f3f4f6;
}

.table-section tr.zero-stock td {
  color: #dc2626;
  font-weight: bold;
}
</style>

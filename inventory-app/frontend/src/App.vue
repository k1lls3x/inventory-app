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
            <LineChart
              v-if="weeklyStockChartData.datasets[0].data.length"
              :data="weeklyStockChartData"
            />
          </div>

          <div class="chart-card">
            <p class="title">Оборот по складам</p>
            <LineChart
              v-if="turnoverLineChartData.datasets[0].data.length"
              :data="turnoverLineChartData"
            />
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

      <section v-else-if="currentTab === 'Остатки'" class="table-section">
        <div class="filter-controls">
          <div class="filter-row">
            <div class="filter-group">
              <label for="warehouse">📦 Склад</label>
              <select v-model="selectedWarehouseId" class="input">
                <option value="0">Все склады</option>
                <option
                  v-for="wh in warehouses"
                  :key="wh.warehouse_id"
                  :value="wh.warehouse_id"
                >
                  {{ wh.name }}
                </option>
              </select>
            </div>

            <div class="filter-group">
              <label for="search">🔍 Поиск</label>
              <input
                type="text"
                class="input"
                v-model="searchQuery"
                placeholder="Название, SKU или склад"
              />
            </div>

            <div class="filter-group button-group">
              <label>&nbsp;</label>
              <button class="add-button" @click="openAddModal">➕ Добавить остаток</button>
            </div>
          </div>
        </div>

        <div
  v-if="showAddModal"
  class="modal-overlay"

  @click.self="closeAddModal"
>
  <div class="modal">
    <h3>Добавить остаток</h3>

    <div class="form-group">
      <label for="item">Товар</label>
      <select v-model.number="newStock.item_id">
        <option disabled value="0">Выберите товар</option>
        <option
          v-for="item in items"
          :key="item.item_id"
          :value="item.item_id"
        >
          {{ item.name }} ({{ item.sku }})
        </option>
      </select>
    </div>

    <div class="form-group">
      <label for="warehouse">Склад</label>
      <select v-model.number="newStock.warehouse_id">
        <option
          v-for="wh in warehouses"
          :key="wh.warehouse_id"
          :value="wh.warehouse_id"
        >
          {{ wh.name }}
        </option>
      </select>
    </div>

    <div class="form-group">
      <label for="quantity">Количество</label>
      <input type="number" v-model.number="newStock.quantity" min="1" />
    </div>

    <div class="modal-actions">
      <button @click="confirmAddStock">💾 Сохранить</button>
      <button @click="closeAddModal">❌ Отмена</button>
    </div>
  </div>
</div>
<!-- Модальное окно редактирования остатка -->
<div
  v-if="showEditModal"
  class="modal-overlay"
  @click.self="closeEditModal"
>
  <div class="modal">
    <h3>Редактировать остаток</h3>

    <div class="form-group">
      <label for="item">Товар</label>
      <input type="text" :value="stockToEdit.name" disabled />
    </div>

    <div class="form-group">
      <label for="warehouse">Склад</label>
      <input type="text" :value="stockToEdit.warehouse" disabled />
    </div>

    <div class="form-group">
      <label for="quantity">Количество</label>
      <input type="number" v-model.number="stockToEdit.quantity" min="1" />
    </div>

    <div class="modal-actions">
      <button @click="confirmEditStock">💾 Сохранить</button>
      <button @click="closeEditModal">❌ Отмена</button>
    </div>
  </div>
</div>
        <div class="fade-in">
          <div class="table-header">
    <p class="title">Остатки на складе</p>
    <button class="export-button" @click="exportToExcel">📤 Экспорт в Excel</button>
  </div>
          <div class="chart-segment">
            <BarChart
              v-if="filteredChartData.datasets[0].data.length"
              :data="filteredChartData"
          />
        </div>
          <table>
            <thead>
              <tr>
                <th>Наименование</th>
                <th>Номер</th>
                <th>Склад</th>
                <th>Количество</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="stock in filteredStockList" :key="stock.id">
                <td>{{ stock.name }}</td>
                <td>{{ stock.sku }}</td>
                <td>{{ stock.warehouse }}</td>
                <td>{{ stock.quantity }}</td>
                <td>
                  <button @click="openEditModal(stock)">✏️</button>
                  <button @click="deleteStock(stock.id)">🗑️</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section v-else>
        <p>Раздел "{{ currentTab }}" в разработке...</p>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import BarChart from './components/BarChart.vue'
import LineChart from './components/LineChart.vue'
import { GetWeeklyStockTrend } from '../wailsjs/go/app/App'
import { ChangeStock } from '../wailsjs/go/app/App'
import { GetStockDetails } from '../wailsjs/go/app/App'
import {
  GetDashboard,
  GetTopItems,
  GetTurnoverByWarehouse,
  FindStockByWarehouse,
  GetWarehouses,
  AddStock,
  GetAllItems
} from '../wailsjs/go/app/App'

const tabs = [
  'Дашборд',
  'Остатки',
  'Поставки',
  'Товары',
  'Склады',
  'Поставщики',
  'Движения'
]
const showEditModal = ref(false)
const stockToEdit = ref(null)
const warehouses = ref([])
const selectedWarehouseId = ref(0)
const searchQuery = ref('')
const currentTab = ref('Дашборд')
const stockList = ref([])
const totalStock = ref(0)
const itemCount = ref(0)
const monthlyOrders = ref(0)
const newItems = ref(0)
const topItems = ref([])
const turnoverData = ref([])
const weeklyStockData = ref([])

const weeklyStockChartData = computed(() => ({
  labels: weeklyStockData.value.map(d => formatDate(d.date)),
  datasets: [
    {
      label: 'Остатки',
      data: weeklyStockData.value.map(d => d.total),
      backgroundColor: 'rgba(0, 0, 0, 0.2)',   // полупрозрачный чёрный
      borderColor: '#000',                    // чёрная линия
      pointBackgroundColor: '#000',           // чёрные точки
      pointRadius: (ctx) => ctx.dataIndex === weeklyStockData.value.length - 1 ? 6 : 0,
      pointHoverRadius: 6,
      borderWidth: 3,
      fill: true,
      tension: 0.3
    }
  ]
}))

const filteredChartData = computed(() => {
  return {
    labels: filteredStockList.value.map(item => item.name),
    datasets: [
      {
        label: 'Остатки',
        data: filteredStockList.value.map(item => item.quantity),
        backgroundColor: '#000',
        barThickness: 20
      }
    ]
  }
})

const turnoverBarChartData = computed(() => ({
  labels: turnoverData.value.map(item => item.name),
  datasets: [
    {
      label: 'Поступило',
      data: turnoverData.value.map(item => item.received),
      backgroundColor: '#000',
      barThickness: 20
    }
  ]
}))

const turnoverLineChartData = computed(() => ({
  labels: turnoverData.value.map(d => d.name),
  datasets: [
    {
      label: 'Оборот',
      data: turnoverData.value.map(d => d.received),
      borderColor: '#000',
      tension: 0.4,
      fill: false
    }
  ]
}))

function formatDate(dateStr) {
  const options = { day: '2-digit', month: 'short' }
  return new Date(dateStr).toLocaleDateString('ru-RU', options)
}

function openEditModal(stock) {
  const warehouse = warehouses.value.find(w => w.name === stock.warehouse)

  stockToEdit.value = {
    ...stock,
    item_id: stock.item_id, // или stock.item_id — зависит от названия
    warehouse_id: warehouse?.warehouse_id
  }

  showEditModal.value = true
}

function confirmEditStock() {
  if (!stockToEdit.value || stockToEdit.value.quantity < 0) {
    alert("Количество должно быть больше 0")
    return
  }

  // item_id и warehouse_id обязательно нужны
  ChangeStock(
  stockToEdit.value.item_id,
  stockToEdit.value.warehouse_id,
  stockToEdit.value.quantity
)
    .then(() => {
      closeEditModal()
      // Обновление данных
      if (selectedWarehouseId.value === 0) {
                GetStockDetails().then(data => {
          stockList.value = data.map(s => ({
            id: s.item_id,               // обязательно для v-for :key
            item_id: s.item_id,          // нужно для ChangeStock
            warehouse_id: s.warehouse_id, // нужно для ChangeStock
            name: s.name,
            sku: s.sku,
            warehouse: s.warehouse,
            quantity: s.quantity
          }))
        })      } else {
        FindStockByWarehouse(selectedWarehouseId.value).then(data => {
          stockList.value = data.map(s => ({
            id: s.item_id,
            name: s.name,
            sku: s.sku,
            warehouse: warehouses.value.find(w => w.warehouse_id === selectedWarehouseId.value)?.name || '',
            quantity: s.quantity
          }))
        })
      }
    })
    .catch(err => {
      alert("Ошибка при обновлении остатка")
      console.error(err)
    })
}

const showAddModal = ref(false)
const newStock = ref({ item_id: 0, warehouse_id: 0, quantity: 0 })
const items = ref([])

function openAddModal() {
  showAddModal.value = true
}

function closeAddModal() {
  showAddModal.value = false
  newStock.value = { item_id: 0, warehouse_id: 0, quantity: 0 }
}

function confirmAddStock() {
  const filled = Object.values(newStock.value).every(v => v !== '' && v !== 0)
  if (!filled) {
    alert('Пожалуйста, заполните все поля')
    return
  }

  AddStock(
    newStock.value.item_id,
    newStock.value.quantity,
    newStock.value.warehouse_id
  ).then(() => {
    closeAddModal()
    if (selectedWarehouseId.value === 0) {
      GetStockDetails().then(data => {
          stockList.value = data.map(s => ({
            id: s.item_id,                // ← важно для v-for :key
            item_id: s.item_id,           // ← нужно для ChangeStock
            warehouse_id: s.warehouse_id, // ← нужно для ChangeStock
            name: s.name,
            sku: s.sku,
            warehouse: s.warehouse,
            quantity: s.quantity
          }))
        })
    } else {
      FindStockByWarehouse(selectedWarehouseId.value).then(data => {
        stockList.value = data.map(s => ({
          id: s.item_id,
          item_id: s.item_id,
          warehouse_id: s.warehouse_id,
          name: s.name,
          sku: s.sku,
          warehouse: s.warehouse, // ← теперь name уже приходит из бэка
          quantity: s.quantity
        }))
      })
    }
  }).catch(err => {
    alert("Ошибка при добавлении остатка")
    console.error(err)
  })
}

function closeEditModal() {
  showEditModal.value = false
  stockToEdit.value = null
}

onMounted(() => {
  GetWeeklyStockTrend().then(data => weeklyStockData.value = data)
  GetAllItems().then(data => items.value = data)
  GetDashboard().then(data => {
    totalStock.value = data.total_stock
    itemCount.value = data.item_count
    monthlyOrders.value = data.monthly_orders
    newItems.value = data.new_items
  })
  GetStockDetails().then(data => {
    stockList.value = data.map(s => ({
      id: s.item_id,               // обязательно для :key
      item_id: s.item_id,          // нужно для ChangeStock
      warehouse_id: s.warehouse_id,
      name: s.name,
      sku: s.sku,
      warehouse: s.warehouse,
      quantity: s.quantity
    }))
  })
  GetTopItems().then(data => topItems.value = data)
  GetWarehouses().then(data => warehouses.value.push(...data))
  GetTurnoverByWarehouse().then(data => turnoverData.value = data)
})

watch(currentTab, (tab) => {
  if (tab === 'Дашборд') {
    GetDashboard().then(data => {
      totalStock.value = data.total_stock
      itemCount.value = data.item_count
      monthlyOrders.value = data.monthly_orders
      newItems.value = data.new_items
    })
    GetTopItems().then(data => topItems.value = data)
    GetTurnoverByWarehouse().then(data => turnoverData.value = data)
  }
})

watch(selectedWarehouseId, (id) => {
  const warehouseId = Number(id)
  if (warehouseId === 0) {
    GetStockDetails().then(data => stockList.value = data)
  } else {
    FindStockByWarehouse(warehouseId).then(data => {
      stockList.value = data.map(s => ({
        id: s.item_id,
        name: s.name,
        sku: s.sku,
        warehouse: warehouses.value.find(w => w.warehouse_id === warehouseId)?.name || '',
        quantity: s.quantity
      }))
    })
  }
})

const filteredStockList = computed(() =>
  stockList.value.filter(item =>
    item.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    item.sku.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    item.warehouse.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
)
</script>

<style scoped>
.dashboard {
  padding: 2rem;
  background-color: #f9f9f9;
}
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
}

.modal {
  background-color: white;
  padding: 1.5rem;
  border-radius: 12px;
  width: 360px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.2);
}

.modal h3 {
  margin-bottom: 1rem;
  font-size: 1.2rem;
}

.form-group {
  margin-bottom: 1rem;
  display: flex;
  flex-direction: column;
}

.form-group input,
.form-group select {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 6px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

.chart-segment {
  margin-bottom: 2rem;
}

.modal-actions button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: bold;
}

.modal-actions button:first-child {
  background-color: #2563eb;
  color: white;
}

.modal-actions button:last-child {
  background-color: #e5e7eb;
}
.stock-filters {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1rem;
  max-width: 500px;
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

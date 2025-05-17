<template>
    <LoginForm v-if="!loggedIn" @login-success="onLoginSuccess" />
   <div v-else class="dashboard">
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
              <button @click="logout">Выйти</button>

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
                  <button @click="deleteStock(stock)">🗑️</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
      <section v-else-if="currentTab === 'Поставки'" class="table-section">

<!-- Фильтры -->
<div class="filter-controls">
  <div class="filter-row">
    <div class="filter-group">
      <label>📅 Дата</label>
      <input
        type="date"
        class="input"
        v-model="selectedDeliveryDate"
        :max="'3030-12-31'" />
    </div>
    <div class="filter-group">
      <label>🔍 Поиск</label>
      <input
        type="text"
        class="input"
        v-model="deliverySearchQuery"
        placeholder="Название, SKU или поставщик">
    </div>
    <div class="filter-group button-group">
      <label>&nbsp;</label>
      <button class="add-button" @click="openAddDeliveryModal">
        ➕ Добавить поставку
      </button>
    </div>
  </div>
</div>

<!-- Модалка ДОБАВИТЬ поставку -->
<div
  v-if="showAddDeliveryModal"
  class="modal-overlay"
  @click.self="closeAddDeliveryModal"
>
  <div class="modal">
    <h3>Добавить поставку</h3>
    <div class="form-group">
      <label for="inbound-item">Товар</label>
      <select v-model.number="newInbound.item_id" id="inbound-item">
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
      <label for="inbound-supplier">Поставщик</label>
      <select v-model.number="newInbound.supplier_id" id="inbound-supplier">
        <option disabled value="0">Выберите поставщика</option>
        <option
          v-for="sup in suppliers"
          :key="sup.supplier_id"
          :value="sup.supplier_id"
        >
          {{ sup.name }}
        </option>
      </select>
    </div>
    <div class="form-group">
      <label for="inbound-warehouse">Склад</label>
      <select v-model.number="newInbound.warehouse_id" id="inbound-warehouse">
        <option disabled value="0">Выберите склад</option>
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
      <label for="inbound-quantity">Количество</label>
      <input
        type="number"
        min="1"
        v-model.number="newInbound.quantity"
        id="inbound-quantity"
      />
    </div>
    <div class="form-group">
      <label for="inbound-date">Дата поступления</label>
      <input
        type="date"
        v-model="newInbound.received_at"
        id="inbound-date"
        :max="'3030-12-31'"
      />
    </div>
    <div class="modal-actions">
      <button @click="confirmAddDelivery">💾 Сохранить</button>
      <button @click="closeAddDeliveryModal">❌ Отмена</button>
    </div>
  </div>
</div>

<!-- Модалка РЕДАКТИРОВАТЬ поставку -->
<div
  v-if="showEditDeliveryModal && deliveryToEdit"
  class="modal-overlay"
  @click.self="closeEditDeliveryModal"
>
  <div class="modal">
    <h3>Редактировать поставку</h3>
    <div class="form-group">
      <label for="edit-inbound-item">Товар</label>
      <select v-model.number="deliveryToEdit.item_id" id="edit-inbound-item">
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
      <label for="edit-inbound-supplier">Поставщик</label>
      <select v-model.number="deliveryToEdit.supplier_id" id="edit-inbound-supplier">
        <option disabled value="0">Выберите поставщика</option>
        <option
          v-for="sup in suppliers"
          :key="sup.supplier_id"
          :value="sup.supplier_id"
        >
          {{ sup.name }}
        </option>
      </select>
    </div>
    <div class="form-group">
      <label for="edit-inbound-warehouse">Склад</label>
      <select v-model.number="deliveryToEdit.warehouse_id" id="edit-inbound-warehouse">
        <option disabled value="0">Выберите склад</option>
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
      <label for="edit-inbound-quantity">Количество</label>
      <input
        type="number"
        min="1"
        v-model.number="deliveryToEdit.quantity"
        id="edit-inbound-quantity"
      />
    </div>
    <div class="form-group">
      <label for="edit-inbound-date">Дата поступления</label>
      <input
        type="date"
        v-model="deliveryToEdit.received_at"
        id="edit-inbound-date"
        :max="'3030-12-31'"
      />
    </div>
    <div class="modal-actions">
      <button @click="confirmEditDelivery">💾 Сохранить</button>
      <button @click="closeEditDeliveryModal">❌ Отмена</button>
    </div>
  </div>
</div>

<!-- Таблица и остальные элементы -->
<div class="fade-in">
  <div class="table-header">
    <p class="title">Поставки</p>
    <button class="export-button" @click="exportDeliveriesToExcel">
      📤 Экспорт в Excel
    </button>
  </div>
  <div class="chart-segment">
    <BarChart
      v-if="filteredDeliveriesChartData.datasets[0].data.length"
      :data="filteredDeliveriesChartData"
    />
  </div>
  <table>
    <thead>
      <tr>
        <th>Дата</th>
        <th>Наименование</th>
        <th>SKU</th>
        <th>Склад</th>
        <th>Поставщик</th>
        <th>Количество</th>
        <th>Действия</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="d in filteredDeliveriesList" :key="d.id">
        <td>{{ formatDate(d.date) }}</td>
        <td>{{ d.name }}</td>
        <td>{{ d.sku }}</td>
        <td>{{ d.warehouse }}</td>
        <td>{{ d.supplier }}</td>
        <td>{{ d.quantity }}</td>
        <td>
          <button @click="openEditDeliveryModal(d)">✏️</button>
          <button @click="deleteDelivery(d)">🗑️</button>
        </td>
      </tr>
    </tbody>
    <div v-if="filteredDeliveriesList.length === 0" class="empty-message" style="text-align:center; color:#888; margin: 1.5rem 0;">
      Нет поставок за выбранную дату
    </div>
  </table>
</div>
</section>

<section v-else-if="currentTab === 'Товары'" class="table-section">
  <div class="filter-controls">
    <div class="filter-row">
      <div class="filter-group">
        <label>🔍 Поиск</label>
        <input
          type="text"
          class="input"
          v-model="itemSearch"
          placeholder="Название, SKU или категория"
        />
      </div>
      <div class="filter-group">
        <label>Категория</label>
        <select v-model="selectedCategory" class="input">
          <option value="">Все категории</option>
          <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
        </select>
      </div>
      <div class="filter-group button-group">
        <label>&nbsp;</label>
        <button class="add-button" @click="openAddItemModal">
          ➕ Добавить товар
        </button>
      </div>
    </div>
  </div>

  <!-- Статистика по товарам -->
  <div class="cards" style="margin-bottom: 1.5rem;">
    <div class="card">
      <p class="title">Всего товаров</p>
      <p class="value">{{ items.length }}</p>
    </div>
    <div class="card">
      <p class="title">Категорий</p>
      <p class="value">{{ categories.length }}</p>
    </div>
    <div class="card">
      <p class="title">Наименьший остаток</p>
      <p class="value" :class="{'note': true, 'positive': minStock > 10, 'negative': minStock <= 10}">
        {{ minStock }}
      </p>
    </div>
    <div class="card">
      <p class="title">Наибольший остаток</p>
      <p class="value">{{ maxStock }}</p>
    </div>
  </div>

  <!-- Таблица товаров -->
  <div class="fade-in">
    <div class="table-header">
      <p class="title">Товары</p>
      <button class="export-button" @click="exportItemsToExcel">
        📤 Экспорт в Excel
      </button>
    </div>
    <table>
      <thead>
        <tr>
          <th>Наименование</th>
          <th>SKU</th>
          <th>Категория</th>
          <th>Ед. изм.</th>
          <th>Мин. остаток</th>
          <th>Описание</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in filteredItems" :key="item.item_id">
          <td>{{ item.name }}</td>
          <td>{{ item.sku }}</td>
          <td>{{ item.category }}</td>
          <td>{{ item.unit }}</td>
          <td>{{ item.min_stock }}</td>
          <td>{{ item.description }}</td>
          <td>
            <button @click="openEditItemModal(item)">✏️</button>
            <button @click="deleteItem(item)">🗑️</button>
          </td>
        </tr>
      </tbody>
      <div v-if="filteredItems.length === 0" class="empty-message" style="text-align:center; color:#888; margin: 1.5rem 0;">
        Нет товаров по фильтру
      </div>
    </table>
  </div>

  <!-- Модалки добавить/редактировать (аналогично поставкам) -->
  <!-- ... -->
</section>
<section v-else>
  <p>Раздел "{{ currentTab }}" в разработке...</p>
</section>
    </main>
  </div>
</div>
</template>
<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import BarChart from './components/BarChart.vue'
import LineChart from './components/LineChart.vue'
import { GetWeeklyStockTrend } from '../wailsjs/go/app/App'
import { ChangeStock } from '../wailsjs/go/app/App'
import { GetStockDetails } from '../wailsjs/go/app/App'
import { RemoveStock } from '../wailsjs/go/app/App'
import { ExportStockToExcel } from '../wailsjs/go/app/App'
import LoginForm from './components/LoginForm.vue'

import {
  GetDashboard,
  GetTopItems,
  GetTurnoverByWarehouse,
  FindStockByWarehouse,
  GetWarehouses,
  AddStock,
  GetAllItems,
  GetInboundDetails,
  GetInboundDetailsByDate,
  AddInbound,
  GetSuppliers,
  DeleteInbound,
  EditInbound

} from '../wailsjs/go/app/App'
const loggedIn = ref(localStorage.getItem('loggedIn') === 'true')
const emit = defineEmits(['login-success'])

function onLoginSuccess() {
  loggedIn.value = true
}
function logout() {
  localStorage.removeItem('loggedIn')
  loggedIn.value = false
}
function handleLogin() {
  error.value = ''
  loading.value = true
  setTimeout(() => {
    if (login.value === 'admin' && password.value === '1234') {
      localStorage.setItem('loggedIn', 'true')
      emit('login-success')
    } else {
      error.value = 'Неверный логин или пароль'
    }
    loading.value = false
  }, 700)
}
const tabs = [
  'Дашборд',
  'Остатки',
  'Поставки',
  'Товары',
  'Склады',
  'Поставщики',
  'Движения'
]

const showEditDeliveryModal = ref(false)
const deliveryToEdit = ref(null)
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
const showAddDeliveryModal = ref(false)
const newInbound = ref({
  item_id: 0,
  supplier_id: 0,
  warehouse_id: 0,
  quantity: 1,
  received_at: "",
})
const suppliers = ref([]) // список поставщиков
function exportToExcel() {
  window.go.app.App.ExportStockToExcel().then(base64data => {
    const binary = atob(base64data);
    const len = binary.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    const blob = new Blob([bytes], { type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = "stock_report.xlsx";
    link.click();
    setTimeout(() => URL.revokeObjectURL(link.href), 1000);
  }).catch(err => {
    alert("Ошибка экспорта: " + err);
  });
}
function openAddDeliveryModal() {
  showAddDeliveryModal.value = true
}

  function openEditDeliveryModal(delivery) {
  let date = delivery.received_at
  if (typeof date === "string" && date.includes(".")) {
    // DD.MM.YYYY -> YYYY-MM-DD
    const [dd, mm, yyyy] = date.split(".");
    date = `${yyyy}-${mm}-${dd}`;
  }
  deliveryToEdit.value = {
    ...delivery,
    received_at: date ? date.substring(0, 10) : ""
  }
  showEditDeliveryModal.value = true
}

// Закрыть модалку
function closeEditDeliveryModal() {
  showEditDeliveryModal.value = false
  deliveryToEdit.value = null
}

// Подтвердить редактирование
function confirmEditDelivery() {
  if (
    !deliveryToEdit.value.item_id ||
    !deliveryToEdit.value.supplier_id ||
    !deliveryToEdit.value.warehouse_id ||
    !deliveryToEdit.value.quantity ||
    deliveryToEdit.value.quantity <= 0 ||
    !deliveryToEdit.value.received_at
  ) {
    alert("Заполните все поля");
    return;
  }
  // Собираем payload для Go backend
  const receivedAt = deliveryToEdit.value.received_at
  ? new Date(deliveryToEdit.value.received_at).toISOString()
  : undefined;

const payload = {
  inbound_id: deliveryToEdit.value.inbound_id,
  item_id: deliveryToEdit.value.item_id,
  supplier_id: deliveryToEdit.value.supplier_id,
  warehouse_id: deliveryToEdit.value.warehouse_id,
  quantity: deliveryToEdit.value.quantity,
  received_at: receivedAt,
  received_by: deliveryToEdit.value.received_by || 1
}
  window.go.app.App.EditInbound(payload).then(() => {
    closeEditDeliveryModal()
    GetInboundDetails().then(data => {
      deliveriesList.value = data || []
    })
  }).catch(err => {
    alert("Ошибка при обновлении поставки")
    console.error(err)
  })
}


function closeAddDeliveryModal() {
  showAddDeliveryModal.value = false
  newInbound.value = {
    item_id: 0,
    supplier_id: 0,
    warehouse_id: 0,
    quantity: 1,
    received_at: "",
  }
}

// Эту функцию вызывай при сохранении поставки
function confirmAddDelivery() {
  // Простая валидация
  if (
    !newInbound.value.item_id ||
    !newInbound.value.supplier_id ||
    !newInbound.value.warehouse_id ||
    !newInbound.value.quantity ||
    newInbound.value.quantity <= 0
  ) {
    alert("Заполните все обязательные поля");
    return;
  }
  const receivedAt = newInbound.value.received_at
  ? new Date(newInbound.value.received_at).toISOString()
  : undefined;
  // Если дата не выбрана, на бэке ставится now()
  const payload = {
    item_id: newInbound.value.item_id,
    supplier_id: newInbound.value.supplier_id,
    warehouse_id: newInbound.value.warehouse_id,
    quantity: newInbound.value.quantity,
    received_at: receivedAt,
  received_by: 1 // <-- id текущего пользователя, если есть логин; пока хардкод
  }
  window.go.app.App.AddInbound(payload).then(() => {
    closeAddDeliveryModal()
    // обновить deliveriesList после добавления
    GetInboundDetails().then(data => {
      deliveriesList.value = data || []
    })
  }).catch(err => {
    alert("Ошибка при добавлении поставки")
    console.error(err)
  })
}

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

function deleteStock(stock) {
  if (!confirm(`Удалить остаток товара "${stock.name}" со склада "${stock.warehouse}"?`)) {
    return
  }

  RemoveStock(stock.stock_id)
    .then(() => {
      const reload = selectedWarehouseId.value === 0
        ? GetStockDetails
        : () => FindStockByWarehouse(selectedWarehouseId.value)

      reload().then(data => {
        stockList.value = data.map(s => ({
          id: s.stock_id,
          stock_id: s.stock_id,
          item_id: s.item_id,
          warehouse_id: s.warehouse_id,
          name: s.name,
          sku: s.sku,
          warehouse: s.warehouse,
          quantity: s.quantity
        }))
      })
    })
    .catch(err => {
      alert("Ошибка при удалении")
      console.error(err)
    })
}


const selectedDeliveryWarehouse = ref(0)
const selectedDeliveryDate = ref("")
const deliverySearchQuery = ref("")
const deliveriesList = ref([])

watch(selectedDeliveryDate, (date) => {
  if (date) {
    GetInboundDetailsByDate(date).then(data => {
      deliveriesList.value = data || [];
    });
  } else {
    GetInboundDetails().then(data => {
      deliveriesList.value = data || [];
    });
  }
});

const filteredDeliveriesList = computed(() =>
  deliveriesList.value.filter(d =>
    d.name.toLowerCase().includes(deliverySearchQuery.value.toLowerCase()) ||
    d.sku.toLowerCase().includes(deliverySearchQuery.value.toLowerCase()) ||
    (d.supplier && d.supplier.toLowerCase().includes(deliverySearchQuery.value.toLowerCase()))
  )
)

const filteredDeliveriesChartData = computed(() => ({
  labels: filteredDeliveriesList.value.map(d => d.name),
  datasets: [
    {
      label: 'Поставки',
      data: filteredDeliveriesList.value.map(d => d.quantity),
      backgroundColor: '#000',
      barThickness: 20
    }
  ]
}))

function deleteDelivery(delivery) {
  console.log("Удаляем поставку:", delivery);
  if (confirm(`Удалить поставку "${delivery.name}" от "${delivery.supplier}"?`)) {
    window.go.app.App.DeleteInbound(delivery.inbound_id || delivery.id)
      .then(() => {
        // обнови deliveriesList
        GetInboundDetails().then(data => {
          deliveriesList.value = data || [];
        });
      })
      .catch(err => {
        alert("Ошибка при удалении поставки");
        console.error(err);
      });
  }
}

function exportDeliveriesToExcel() {
  alert('Заглушка экспорта. Тут будет экспорт в Excel')
}

// Форматирование даты, чтобы не падало

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
  if (!dateStr) return "";
  const d = new Date(dateStr);
  if (isNaN(d)) return dateStr;
  const dd = String(d.getDate()).padStart(2, '0');
  const mm = String(d.getMonth() + 1).padStart(2, '0');
  const yyyy = d.getFullYear();
  return `${dd}.${mm}.${yyyy}`;
}

function openEditModal(stock) {
  const warehouse = warehouses.value.find(w => w.name === stock.warehouse)

  stockToEdit.value = {
    ...stock,
    item_id: stock.item_id,
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
    id: s.stock_id,
    stock_id: s.stock_id,
    item_id: s.item_id,
    warehouse_id: s.warehouse_id,
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
// Категории (генерируются из items)
const categories = computed(() => {
  return Array.from(new Set(items.value.map(i => i.category)));
});

// Поиск и фильтр
const itemSearch = ref('');
const selectedCategory = ref('');

const filteredItems = computed(() =>
  items.value.filter(i =>
    (!selectedCategory.value || i.category === selectedCategory.value) &&
    (
      i.name.toLowerCase().includes(itemSearch.value.toLowerCase()) ||
      i.sku.toLowerCase().includes(itemSearch.value.toLowerCase()) ||
      (i.category && i.category.toLowerCase().includes(itemSearch.value.toLowerCase()))
    )
  )
);

// Минимальный и максимальный остаток
const minStock = computed(() => {
  if (items.value.length === 0) return 0;
  return Math.min(...items.value.map(i => i.min_stock || 0));
});
const maxStock = computed(() => {
  if (items.value.length === 0) return 0;
  return Math.max(...items.value.map(i => i.min_stock || 0));
});

// Заглушки методов
function openAddItemModal() {
  alert('Открыть модалку добавления товара (реализуй сам)');
}
function openEditItemModal(item) {
  alert('Открыть модалку редактирования товара: ' + item.name);
}
function deleteItem(item) {
  if (confirm(`Удалить "${item.name}"?`)) {
    items.value = items.value.filter(i => i.item_id !== item.item_id);
  }
}
function exportItemsToExcel() {
  alert('Заглушка экспорта товаров в Excel');
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
          id: s.stock_id, // ✅ нужно для deleteStock(stock)
          stock_id: s.stock_id, // опционально, если хочешь сохранять явно
          item_id: s.item_id,
          warehouse_id: s.warehouse_id,
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
  GetSuppliers().then(data => suppliers.value = data || [])
    GetStockDetails().then(data => {
    stockList.value = data.map(s => ({
      id: s.stock_id, // 👈 обязательно
      stock_id: s.stock_id,
      item_id: s.item_id,
      warehouse_id: s.warehouse_id,
      name: s.name,
      sku: s.sku,
      warehouse: s.warehouse,
      quantity: s.quantity
    }))
  })
  GetInboundDetails().then(data => {
    deliveriesList.value = data || [];
  }).catch(err => {
    console.error("Ошибка загрузки поставок:", err)
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
    GetStockDetails().then(data => {
      stockList.value = data.map(s => ({
        id: s.stock_id,
        stock_id: s.stock_id,
        item_id: s.item_id,
        warehouse_id: s.warehouse_id,
        name: s.name,
        sku: s.sku,
        warehouse: s.warehouse,
        quantity: s.quantity
      }))
    })
  } else {
    FindStockByWarehouse(warehouseId).then(data => {
      stockList.value = data.map(s => ({
        id: s.stock_id,
        stock_id: s.stock_id,
        item_id: s.item_id,
        warehouse_id: s.warehouse_id,
        name: s.name,
        sku: s.sku,
        warehouse: warehouses.value.find(w => w.warehouse_id === warehouseId)?.name || s.warehouse,
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

</style>

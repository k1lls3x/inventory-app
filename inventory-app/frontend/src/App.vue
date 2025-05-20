<template>
  <LoginForm v-if="!loggedIn" @login-success="onLoginSuccess" />
  <div v-else class="layout">
    <!-- Sidebar -->
    <aside class="sidebar">
  <!-- User Block -->
  <div class="sidebar-user" @click="showProfileModal = true" style="cursor:pointer;">
  <div class="sidebar-user-avatar">
    <svg width="40" height="40" fill="none" viewBox="0 0 40 40">
      <circle cx="20" cy="20" r="20" fill="#e3eaff"/>
      <path d="M20 24c-4 0-7 3-7 7h14c0-4-3-7-7-7Zm0-2a5 5 0 1 0 0-10 5 5 0 0 0 0 10Z" fill="#b7c5ec"/>
    </svg>
  </div>
  <div>
    <div class="sidebar-user-name">{{ user?.full_name || user?.username }}</div>
    <div class="sidebar-user-role">{{ roleName(user.role) }}</div>
  </div>
</div>

  <div class="sidebar-logo">
    <!-- SVG или логотип -->
  </div>
  <nav>
  <button
    v-for="tab in tabs"
    :key="tab"
    :class="{ active: currentTab === tab }"
    @click="currentTab = tab"
  >{{ tab }}</button>
</nav>
  <button class="logout-btn" @click="logout">🚪 Выйти</button>
</aside>

    <!-- Main Content -->
    <div class="main-content">
      <header class="main-header">
        <h1>Складская система</h1>
        <span class="username">Добро пожаловать!</span>
      </header>

      <main>
        <!-- Дашборд -->
        <section v-if="currentTab === 'Дашборд'">
          <div class="cards">
            <div class="card highlight animate-card">
              <p class="title">Всего остатков</p>
              <p class="value">{{ totalStock }}</p>
              <p class="note positive">+15% за месяц</p>
            </div>
            <div class="card animate-card">
              <p class="title">Товаров</p>
              <p class="value">{{ itemCount }}</p>
              <p class="note" v-if="newItems > 0">+{{ newItems }} новых за месяц</p>
            </div>
            <div class="card animate-card">
              <p class="title">Поставки</p>
              <p class="value">{{ monthlyOrders }}</p>
              <p class="note">в этом месяце</p>
            </div>
          </div>
          <div class="charts-table-wrap">
            <div class="chart-card animate-chart">
              <p class="title">Остатки за неделю</p>
              <LineChart v-if="weeklyStockChartData.datasets[0].data.length" :data="weeklyStockChartData" />
            </div>
            <div class="chart-card animate-chart">
              <p class="title">Оборот по складам</p>
              <LineChart v-if="turnoverLineChartData.datasets[0].data.length" :data="turnoverLineChartData" />
            </div>
          </div>
          <div class="table-section">
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
            <div v-if="topItems.length === 0" class="empty-message">
              Нет данных для отображения
            </div>
          </div>
        </section>
 <!-- Пользователи (видно только админу) -->
 <section v-if="currentTab === 'Пользователи' && user?.role === 'admin'">
    <UserTable />
  </section>
        <!-- Остатки -->
        <section v-else-if="currentTab === 'Остатки'">
          <div class="filters-bar">
            <div class="filter-group">
              <label>📦 Склад</label>
              <select v-model="selectedWarehouseId" class="input">
                <option value="0">Все склады</option>
                <option v-for="wh in warehouses" :key="wh.warehouse_id" :value="wh.warehouse_id">
                  {{ wh.name }}
                </option>
              </select>
            </div>
            <div class="filter-group">
              <label>🔍 Поиск</label>
              <input type="text" class="input" v-model="searchQuery" placeholder="Название, SKU или склад" />
            </div>
            <div class="filter-group button-group">
              <label>&nbsp;</label>
              <button class="add-button" @click="openAddModal">➕ Добавить остаток</button>
            </div>
          </div>

          <!-- Модалка добавления остатка -->
          <div v-if="showAddModal" class="modal-overlay" @click.self="closeAddModal">
            <div class="modal">
              <h3>Добавить остаток</h3>
              <div class="form-group">
                <label for="item">Товар</label>
                <select v-model.number="newStock.item_id">
                  <option disabled value="0">Выберите товар</option>
                  <option v-for="item in items" :key="item.item_id" :value="item.item_id">
                    {{ item.name }} ({{ item.sku }})
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label for="warehouse">Склад</label>
                <select v-model.number="newStock.warehouse_id">
                  <option v-for="wh in warehouses" :key="wh.warehouse_id" :value="wh.warehouse_id">
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

          <!-- Модалка редактирования остатка -->
          <div v-if="showEditModal" class="modal-overlay" @click.self="closeEditModal">
            <div class="modal">
              <h3>Редактировать остаток</h3>
              <div class="form-group">
                <label for="item">Товар</label>
                <input type="text" :value="stockToEdit?.name" disabled />
              </div>
              <div class="form-group">
                <label for="warehouse">Склад</label>
                <input type="text" :value="stockToEdit?.warehouse" disabled />
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

          <div class="charts-table-wrap">
            <div class="chart-card animate-chart">
              <BarChart v-if="filteredChartData.datasets[0].data.length" :data="filteredChartData" />
            </div>
            <div class="table-section animate-table">
              <div class="table-header">
                <p class="title">Остатки на складе</p>
                <button class="export-button" @click="exportToExcel">📤 Экспорт в Excel</button>
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
                      <div class="action-buttons">
                        <button class="action-btn edit" @click="openEditModal(stock)">✏️</button>
                        <button class="action-btn delete" @click="deleteStock(stock)">🗑️</button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
              <div v-if="filteredStockList.length === 0" class="empty-message">
                Нет остатков для отображения
              </div>
            </div>
          </div>
        </section>

        <!-- Поставки -->
        <section v-else-if="currentTab === 'Поставки'">
          <div class="filters-bar">
            <div class="filter-group">
              <label>📅 Дата</label>
              <input type="date" class="input" v-model="selectedDeliveryDate" :max="'3030-12-31'" />
            </div>
            <div class="filter-group">
              <label>🔍 Поиск</label>
              <input type="text" class="input" v-model="deliverySearchQuery" placeholder="Название, SKU или поставщик" />
            </div>
            <div class="filter-group button-group">
              <label>&nbsp;</label>
              <button class="add-button" @click="openAddDeliveryModal">➕ Добавить поставку</button>
            </div>
          </div>

          <!-- Модалка добавления поставки -->
          <div v-if="showAddDeliveryModal" class="modal-overlay" @click.self="closeAddDeliveryModal">
            <div class="modal">
              <h3>Добавить поставку</h3>
              <div class="form-group">
                <label for="inbound-item">Товар</label>
                <select v-model.number="newInbound.item_id" id="inbound-item">
                  <option disabled value="0">Выберите товар</option>
                  <option v-for="item in items" :key="item.item_id" :value="item.item_id">
                    {{ item.name }} ({{ item.sku }})
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label for="inbound-supplier">Поставщик</label>
                <select v-model.number="newInbound.supplier_id" id="inbound-supplier">
                  <option disabled value="0">Выберите поставщика</option>
                  <option v-for="sup in suppliers" :key="sup.supplier_id" :value="sup.supplier_id">
                    {{ sup.name }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label for="inbound-warehouse">Склад</label>
                <select v-model.number="newInbound.warehouse_id" id="inbound-warehouse">
                  <option disabled value="0">Выберите склад</option>
                  <option v-for="wh in warehouses" :key="wh.warehouse_id" :value="wh.warehouse_id">
                    {{ wh.name }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label for="inbound-quantity">Количество</label>
                <input type="number" min="1" v-model.number="newInbound.quantity" id="inbound-quantity" />
              </div>
              <div class="form-group">
                <label for="inbound-date">Дата поступления</label>
                <input type="date" v-model="newInbound.received_at" id="inbound-date" :max="'3030-12-31'" />
              </div>
              <div class="modal-actions">
                <button @click="confirmAddDelivery">💾 Сохранить</button>
                <button @click="closeAddDeliveryModal">❌ Отмена</button>
              </div>
            </div>
          </div>

          <!-- Модалка редактирования поставки -->
          <div v-if="showEditDeliveryModal && deliveryToEdit" class="modal-overlay" @click.self="closeEditDeliveryModal">
            <div class="modal">
              <h3>Редактировать поставку</h3>
              <div class="form-group">
                <label for="edit-inbound-item">Товар</label>
                <select v-model.number="deliveryToEdit.item_id" id="edit-inbound-item">
                  <option disabled value="0">Выберите товар</option>
                  <option v-for="item in items" :key="item.item_id" :value="item.item_id">
                    {{ item.name }} ({{ item.sku }})
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label for="edit-inbound-supplier">Поставщик</label>
                <select v-model.number="deliveryToEdit.supplier_id" id="edit-inbound-supplier">
                  <option disabled value="0">Выберите поставщика</option>
                  <option v-for="sup in suppliers" :key="sup.supplier_id" :value="sup.supplier_id">
                    {{ sup.name }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label for="edit-inbound-warehouse">Склад</label>
                <select v-model.number="deliveryToEdit.warehouse_id" id="edit-inbound-warehouse">
                  <option disabled value="0">Выберите склад</option>
                  <option v-for="wh in warehouses" :key="wh.warehouse_id" :value="wh.warehouse_id">
                    {{ wh.name }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label for="edit-inbound-quantity">Количество</label>
                <input type="number" min="1" v-model.number="deliveryToEdit.quantity" id="edit-inbound-quantity" />
              </div>
              <div class="form-group">
                <label for="edit-inbound-date">Дата поступления</label>
                <input type="date" v-model="deliveryToEdit.received_at" id="edit-inbound-date" :max="'3030-12-31'" />
              </div>
              <div class="modal-actions">
                <button @click="confirmEditDelivery">💾 Сохранить</button>
                <button @click="closeEditDeliveryModal">❌ Отмена</button>
              </div>
            </div>
          </div>

          <div class="charts-table-wrap">
            <div class="chart-card animate-chart">
              <BarChart v-if="filteredDeliveriesChartData.datasets[0].data.length" :data="filteredDeliveriesChartData" />
            </div>
            <div class="table-section animate-table">
              <div class="table-header">
                <p class="title">Поставки</p>
                <button class="export-button" @click="exportDeliveriesToExcel">📤 Экспорт в Excel</button>
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
                      <div class="action-buttons">
                        <button class="action-btn edit" @click="openEditDeliveryModal(d)">✏️</button>
                        <button class="action-btn delete" @click="deleteDelivery(d)">🗑️</button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
              <div v-if="filteredDeliveriesList.length === 0" class="empty-message">
                Нет поставок за выбранную дату
              </div>
            </div>
          </div>
          <div>{{ items.length }}</div>
        </section>

        <!-- Товары -->
        <section v-else-if="currentTab === 'Товары'">
          <div class="filters-bar">
        <div class="filter-group">
          <label>🔍 Поиск</label>
          <input type="text" class="input" v-model="itemSearch" placeholder="Название или SKU" />
        </div>
        <div class="filter-group button-group">
          <label>&nbsp;</label>
          <button class="add-button" @click="openAddItemModal">➕ Добавить товар</button>
        </div>
      </div>
      <div class="cards">
          <div class="card animate-card">
            <p class="title">Всего товаров</p>
            <p class="value">{{ items.length }}</p>
          </div>
          <div class="card animate-card">
            <p class="title">Средняя цена товара</p>
            <p class="value">{{ averagePrice }}</p>
          </div>
          <div class="card animate-card">
            <p class="title">Наибольший остаток</p>
            <p class="value">
              {{ maxActualStock }}
            </p>
          </div>
        </div>

          <div class="table-section animate-table">
            <div class="table-header">
              <p class="title">Товары</p>
              <button class="export-button" @click="exportItemsToExcel">📤 Экспорт в Excel</button>
            </div>
            <table>
                  <thead>
                    <tr>
                      <th>Наименование</th>
                      <th>SKU</th>
                      <th>Ед. изм.</th>
                      <th>Описание</th>
                      <th>Цена</th>
                      <th>Себестоимость</th>
                      <th>Действия</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="item in filteredItems" :key="item.item_id">
                      <td>{{ item.name }}</td>
                      <td>{{ item.sku }}</td>
                      <td>{{ item.uom }}</td>
                      <td>{{ item.description }}</td>
                      <td>
                        {{ item.price != null ? Number(item.price).toLocaleString('ru-RU', { minimumFractionDigits: 2 }) : '—' }}
                      </td>
                      <td>{{ item.cost }}</td>
                      <td>
                        <div class="action-buttons">
                          <button class="action-btn edit" @click="openEditItemModal(item)">✏️</button>
                          <button class="action-btn delete" @click="deleteItem(item)">🗑️</button>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>

            <div v-if="filteredItems.length === 0" class="empty-message">
              Нет товаров по фильтру
            </div>
          </div>
          <div v-if="showAddItemModal" class="modal-overlay" @click.self="showAddItemModal = false">
  <div class="modal">
    <h3>Добавить товар</h3>
    <div class="form-group"><label>SKU</label><input v-model="newItem.sku" /></div>
<div class="form-group"><label>Наименование</label><input v-model="newItem.name" /></div>
<div class="form-group"><label>Описание</label><input v-model="newItem.description" /></div>
<div class="form-group"><label>Ед. изм.</label><input v-model="newItem.uom" /></div>
<div class="form-group"><label>Цена</label><input type="number" v-model.number="newItem.price" /></div>
<div class="form-group"><label>Себестоимость</label><input type="number" v-model.number="newItem.cost" /></div>

    <div class="modal-actions">
      <button @click="confirmAddItem">💾 Сохранить</button>
      <button @click="showAddItemModal = false">❌ Отмена</button>
    </div>
  </div>
</div>
<div v-if="showEditItemModal" class="modal-overlay" @click.self="showEditItemModal = false">
  <div class="modal item-edit-modal">
    <h3>Редактировать товар</h3>
    <form @submit.prevent="confirmEditItem" autocomplete="off">
      <div class="form-group"><label>SKU</label>
        <input v-model="itemToEdit.sku" disabled class="input-modern" />
      </div>
      <div class="form-group"><label>Наименование</label>
        <input v-model="itemToEdit.name" class="input-modern" />
      </div>
      <div class="form-group"><label>Описание</label>
        <textarea v-model="itemToEdit.description" rows="2" class="input-modern" style="resize:vertical; min-height:36px;" />
      </div>
      <div class="form-row">
        <div class="form-group half"><label>Категория</label>
          <input v-model="itemToEdit.category" class="input-modern" />
        </div>
        <div class="form-group half"><label>Ед. изм.</label>
          <input v-model="itemToEdit.uom" class="input-modern" />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group half"><label>Мин. остаток</label>
          <input type="number" v-model.number="itemToEdit.reorder_level" min="0" class="input-modern" />
        </div>
        <div class="form-group half"><label>Партия для дозакупки</label>
          <input type="number" v-model.number="itemToEdit.reorder_qty" min="0" class="input-modern" />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group half"><label>Цена</label>
          <input type="number" v-model.number="itemToEdit.price" min="0" step="0.01" class="input-modern" />
        </div>
        <div class="form-group half"><label>Себестоимость</label>
          <input type="number" v-model.number="itemToEdit.cost" min="0" step="0.01" class="input-modern" />
        </div>
      </div>
      <div class="modal-actions modal-actions-row">
        <button type="submit" class="main-btn-strong">💾 Сохранить</button>
        <button type="button" class="main-btn-ghost" @click="showEditItemModal = false">❌ Отмена</button>
      </div>
    </form>
  </div>
</div>

<!-- Поставщики -->
<section v-else-if="currentTab === 'Поставщики'">


  <div class="filters-bar">
    <div class="filter-group">
      <label>🔍 Поиск</label>
      <input type="text" class="input" v-model="supplierSearch" placeholder="Название или ИНН" />
    </div>
    <div class="filter-group button-group">
      <label>&nbsp;</label>
      <button class="add-button" @click="openAddSupplierModal">➕ Добавить поставщика</button>
    </div>
  </div>

  <div class="cards">
    <div class="card animate-card">
      <p class="title">Всего поставщиков</p>
      <p class="value">{{ suppliers.length }}</p>
    </div>
  </div>

  <div class="table-section animate-table">
    <div class="table-header">
      <p class="title">Поставщики</p>
      <button class="export-button" @click="exportSuppliersToExcel">📤 Экспорт в Excel</button>
    </div>
    <table>
      <thead>
        <tr>
          <th>Название</th>
          <th>ИНН</th>
          <th>Контакт</th>
          <th>Телефон</th>
          <th>Email</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in filteredSuppliers" :key="s.supplier_id">
          <td>{{ s.name }}</td>
          <td>{{ s.inn }}</td>
          <td>{{ s.contact_person }}</td>
          <td>{{ s.phone }}</td>
          <td>{{ s.email }}</td>
          <td>
            <div class="action-buttons">
              <button class="action-btn edit" @click="openEditSupplierModal(s)">✏️</button>
              <button class="action-btn delete" @click="deleteSupplier(s)">🗑️</button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="filteredSuppliers.length === 0" class="empty-message">
      Нет поставщиков по фильтру
    </div>
  </div>

  <!-- Модалка добавления -->
  <div v-if="showAddSupplierModal" class="modal-overlay" @click.self="showAddSupplierModal = false">
    <div class="modal">
      <h3>Добавить поставщика</h3>
      <div class="form-group"><label>Название</label><input v-model="newSupplier.name" /></div>
      <div class="form-group"><label>ИНН</label><input v-model="newSupplier.inn" /></div>
      <div class="form-group"><label>Контакт</label><input v-model="newSupplier.contact_person" /></div>
      <div class="form-group"><label>Телефон</label><input v-model="newSupplier.phone" /></div>
      <div class="form-group"><label>Email</label><input v-model="newSupplier.email" /></div>
      <div class="modal-actions">
        <button @click="confirmAddSupplier">💾 Сохранить</button>
        <button @click="showAddSupplierModal = false">❌ Отмена</button>
      </div>
    </div>
  </div>

  <!-- Модалка редактирования -->
  <div v-if="showEditSupplierModal" class="modal-overlay" @click.self="showEditSupplierModal = false">
    <div class="modal">
      <h3>Редактировать поставщика</h3>
      <div class="form-group"><label>Название</label><input v-model="supplierToEdit.name" /></div>
      <div class="form-group"><label>ИНН</label><input v-model="supplierToEdit.inn" /></div>
      <div class="form-group"><label>Контакт</label><input v-model="supplierToEdit.contact_person" /></div>
      <div class="form-group"><label>Телефон</label><input v-model="supplierToEdit.phone" /></div>
      <div class="form-group"><label>Email</label><input v-model="supplierToEdit.email" /></div>
      <div class="modal-actions">
        <button @click="confirmEditSupplier">💾 Сохранить</button>
        <button @click="showEditSupplierModal = false">❌ Отмена</button>
      </div>
    </div>
  </div>
</section>

        </section>

        <!-- Другое (заглушка) -->
        <section v-else>
          <p>Раздел "{{ currentTab }}" в разработке...</p>
        </section>
      </main>
    </div>
   <!-- МОДАЛКА ПРОФИЛЯ - вставь в свой <template> -->
    <div v-if="showProfileModal" class="modal-overlay" @click.self="showProfileModal = false">
  <div class="modal profile-modal-modern">
    <!-- Заголовок и аватар -->
    <div class="profile-header-modern">
      <div class="profile-avatar-modern accent-avatar">
        <svg width="58" height="58" viewBox="0 0 58 58">
          <defs>
            <linearGradient id="avatar-gradient" x1="0" y1="0" x2="1" y2="1">
              <stop offset="0%" stop-color="#a5b6fa"/>
              <stop offset="100%" stop-color="#2563eb"/>
            </linearGradient>
          </defs>
          <circle cx="29" cy="29" r="29" fill="url(#avatar-gradient)"/>
          <path d="M29 36c-6 0-11 5-11 11h22c0-6-5-11-11-11Zm0-4a7 7 0 1 0 0-14 7 7 0 0 0 0 14Z" fill="#f3f6fd"/>
        </svg>
      </div>
      <div>
        <div class="profile-title-modern main-name">{{ user.full_name }}</div>
        <div class="profile-role-modern">{{ roleName(user.role) }}</div>
      </div>
    </div>

    <!-- Блок с инфой -->
    <div class="profile-info-modern-rich profile-info-compact">
      <div class="info-row-rich">
        <span class="info-icon-circle">
          <svg width="18" height="18" fill="none" viewBox="0 0 18 18"><circle cx="9" cy="9" r="9" fill="#e3eaff"/><path d="M9 12c-2.2 0-4 1.1-4 2.1v.4h8v-.4C13 13.1 11.2 12 9 12zm0-1.1A2.1 2.1 0 1 0 9 6a2.1 2.1 0 0 0 0 4.2z" fill="#2563eb"/></svg>
        </span>
        <span class="info-label-rich">Логин</span>
        <span class="info-value-rich">{{ user.username }}</span>
      </div>
      <div class="info-row-rich">
        <span class="info-icon-circle">
          <svg width="18" height="18" fill="none" viewBox="0 0 18 18"><rect width="18" height="18" rx="6" fill="#e3eaff"/><path d="M5.7 7.2h6.6v1.2H5.7v-1.2zm0 2h6.6v1.2H5.7v-1.2z" fill="#2563eb"/></svg>
        </span>
        <span class="info-label-rich">Роль</span>
        <span class="info-value-rich">{{ roleName(user.role) }}</span>
      </div>
    </div>

    <!-- Блок смены пароля -->
    <div class="profile-info-modern-rich profile-password-compact">
      <div class="profile-change-title-modern">Смена пароля</div>
      <div class="profile-change-fields-modern profile-fields-spaced">
        <input
          type="password"
          v-model="oldPassword"
          placeholder="Старый пароль"
          class="input-modern input-shadow"
        />
        <input
          type="password"
          v-model="newPassword"
          placeholder="Новый пароль"
          class="input-modern input-shadow"
        />
        <input
          type="password"
          v-model="repeatPassword"
          placeholder="Повторите новый пароль"
          class="input-modern input-shadow"
        />
      </div>
      <transition name="fade">
        <div v-if="profileError" class="error-msg-modern">{{ profileError }}</div>
      </transition>
      <transition name="fade">
        <div v-if="profileSuccess" class="success-msg-modern">{{ profileSuccess }}</div>
      </transition>
    </div>

    <!-- Кнопки в самом низу -->
    <div class="profile-actions-modern buttons-bottom profile-actions-outside">
      <button @click="changePassword" class="change-btn-modern main-btn-strong">Сменить пароль</button>
      <button @click="showProfileModal = false" class="close-btn-modern main-btn-ghost">Закрыть</button>
    </div>
  </div>
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
  EditInbound,
  ChangePassword,
  GetItems,
  AddItem,
  UpdateItem,
  RemoveItem

} from '../wailsjs/go/app/App'
const loggedIn = ref(localStorage.getItem('loggedIn') === 'true')
const emit = defineEmits(['login-success'])

function onLoginSuccess(userData) {
  user.value = userData        // <-- вот здесь присваивай!
  loggedIn.value = true
  localStorage.setItem('loggedIn', 'true')
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
const user = ref(null) // по умолчанию null
// Подгружай реального пользователя с бэка после логина
function roleName(role) {
  switch (role) {
    case 'admin': return 'Администратор'
    case 'manager': return 'Менеджер'
    case 'worker': return 'Сотрудник'
    default: return 'Пользователь'
  }
}
const showProfileModal = ref(false)
const oldPassword = ref('')
const newPassword = ref('')
const repeatPassword = ref('')
const profileError = ref('')
const profileSuccess = ref('')
const tabs = computed(() => {
  if (!user.value) return [];
  if (user.value.role === 'admin') {
    return [
      'Дашборд',
      'Остатки',
      'Поставки',
      'Товары',
      'Склады',
      'Поставщики',
      'Движения',
      'Пользователи' // Только для админа!
    ];
  }
  // Для менеджера
  if (user.value.role === 'manager') {
    return [
      'Дашборд',
      'Остатки',
      'Поставки',
      'Товары',
      'Склады',
      'Поставщики',
      'Движения'
    ];
  }
  // Для сотрудника
  return [
    'Дашборд',
    'Остатки',
    'Поставки',
    'Движения'
  ];
});

async function changePassword() {
  profileError.value = ''
  profileSuccess.value = ''
  if (!oldPassword.value || !newPassword.value || !repeatPassword.value) {
    profileError.value = 'Заполните все поля'
    return
  }
  if (newPassword.value !== repeatPassword.value) {
    profileError.value = 'Пароли не совпадают'
    return
  }
  try {
    // используем username (логин)
    await window.go.app.App.ChangePassword(user.value.username, oldPassword.value, newPassword.value)
    profileSuccess.value = 'Пароль успешно изменён'
    oldPassword.value = newPassword.value = repeatPassword.value = ''
  } catch (e) {
    profileError.value = e?.message || 'Ошибка смены пароля'
  }
}
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
const suppliers = ref([]);
const supplierSearch = ref('');
const filteredSuppliers = computed(() =>
  suppliers.value.filter(s =>
    (s.name || '').toLowerCase().includes(supplierSearch.value.toLowerCase()) ||
    (s.inn || '').toLowerCase().includes(supplierSearch.value.toLowerCase())
  )
);
const showAddSupplierModal = ref(false);
const showEditSupplierModal = ref(false);
const newSupplier = ref({ name: '', inn: '', contact_person: '', phone: '', email: '' });
const supplierToEdit = ref({});

function openAddSupplierModal() { showAddSupplierModal.value = true }
function confirmAddSupplier() {
  suppliers.value.push({ ...newSupplier.value, supplier_id: Date.now() })
  showAddSupplierModal.value = false
  newSupplier.value = { name: '', inn: '', contact_person: '', phone: '', email: '' }
}
function openEditSupplierModal(s) { supplierToEdit.value = { ...s }; showEditSupplierModal.value = true }
function confirmEditSupplier() {
  const idx = suppliers.value.findIndex(x => x.supplier_id === supplierToEdit.value.supplier_id)
  if (idx !== -1) suppliers.value[idx] = { ...supplierToEdit.value }
  showEditSupplierModal.value = false
}
function deleteSupplier(s) {
  suppliers.value = suppliers.value.filter(x => x.supplier_id !== s.supplier_id)
}
function exportSuppliersToExcel() {
  alert('Заглушка экспорта поставщиков')
}

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
const averagePrice = computed(() => {
  if (!items.value.length) return '—'
  // Игнорируем товары без цены (null или 0 можно убрать по желанию)
  const filtered = items.value.filter(i => i.price !== null && i.price !== undefined)
  if (!filtered.length) return '—'
  const sum = filtered.reduce((acc, i) => acc + Number(i.price), 0)
  return (sum / filtered.length).toLocaleString('ru-RU', { minimumFractionDigits: 2 })
})

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
    i.name.toLowerCase().includes(itemSearch.value.toLowerCase()) ||
    i.sku.toLowerCase().includes(itemSearch.value.toLowerCase())
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
const minActualStock = computed(() => {
  if (!totalStockPerItem.value.length) return 0
  return Math.min(...totalStockPerItem.value)
})
const maxActualStock = computed(() => {
  if (!totalStockPerItem.value.length) return 0
  return Math.max(...totalStockPerItem.value)
})


const totalStockPerItem = computed(() => {
  const map = {}
  for (const s of stockList.value) {
    map[s.item_id] = (map[s.item_id] || 0) + s.quantity
  }
  return Object.values(map)
})

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
const showAddItemModal = ref(false)
const showEditItemModal = ref(false)
const itemToEdit = ref(null)

const newItem = ref({
  sku: "",
  name: "",
  description: "",
  uom: "",
  reorder_level: 0,
  reorder_qty: 0,
  price: 0,
  cost: 0,
  category: ""
})
function openAddItemModal() {
  Object.assign(newItem.value, {
    sku: "",
    name: "",
    description: "",
    uom: "",
    reorder_level: 0,
    reorder_qty: 0,
    price: 0,
    cost: 0,
    category: ""
  })
  showAddItemModal.value = true
}

function openEditItemModal(item) {
  itemToEdit.value = { ...item }
  showEditItemModal.value = true
}
async function confirmAddItem() {
  // простая валидация
  if (!newItem.value.sku || !newItem.value.name) {
    alert("Заполните все обязательные поля (артикул и наименование)");
    return;
  }
  try {
    await AddItem(newItem.value)
    showAddItemModal.value = false
    // обнови список
    items.value = await GetItems() || []
  } catch (e) {
    alert('Ошибка при добавлении: ' + (e?.message || ''))
  }
}

async function confirmEditItem() {
  if (!itemToEdit.value.sku || !itemToEdit.value.name) {
    alert("Заполните обязательные поля");
    return;
  }
  // Явно ставим null, если поля пустые (или 0? зависит от бизнес-логики)
  if (itemToEdit.value.price === "") itemToEdit.value.price = null;
  if (itemToEdit.value.cost === "") itemToEdit.value.cost = null;
  try {
    await UpdateItem(itemToEdit.value)
    showEditItemModal.value = false
    items.value = await GetItems() || []
  } catch (e) {
    alert('Ошибка при обновлении: ' + (e?.message || ''))
  }
}
~

async function deleteItem(item) {
  if (!confirm(`Удалить товар "${item.name}"?`)) return
  try {
    await RemoveItem(item.sku)
    items.value = await GetItems() || []
  } catch (e) {
    alert('Ошибка при удалении: ' + (e?.message || ''))
  }
}

function closeEditModal() {
  showEditModal.value = false
  stockToEdit.value = null
}

onMounted(async () => {
  if (loggedIn.value && !user.value) {
    try {
      const userData = await window.go.app.App.GetCurrentUser();
      user.value = userData;
    } catch (err) {
      loggedIn.value = false;
      localStorage.removeItem('loggedIn');
      user.value = null;
    }
  }
  console.log('items:', items.value)
  GetItems().then(data => items.value = data || []);
  GetWeeklyStockTrend().then(data => weeklyStockData.value = data)
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
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 100;
  background: rgba(55, 65, 81, 0.23);
  display: flex;
  align-items: center;
  justify-content: center;
}

.profile-modal-modern {
  min-height: 545px;
  width: 410px;
  max-width: 94vw;
  background: #f9fbff;
  border-radius: 20px;
  box-shadow: 0 10px 48px 0 #20223622;
  padding: 2.7rem 2.2rem 2rem 2.2rem;
  display: flex;
  flex-direction: column;
  gap: 1.8rem;
  animation: modalPopIn .23s cubic-bezier(.44,1.6,.41,.98);
  position: relative;
}

@keyframes modalPopIn {
  0% { transform: scale(0.96) translateY(20px); opacity: 0;}
  100% { transform: scale(1) translateY(0); opacity: 1;}
}

.profile-header-modern {
  display: flex;
  align-items: center;
  gap: 1.1rem;
}
.profile-avatar-modern {
  width: 58px;
  height: 58px;
  background: linear-gradient(120deg, #e3eaff 70%, #f3f8ff 100%);
  border-radius: 50%;
  box-shadow: 0 3px 18px #e2eaff51;
  display: flex;
  align-items: center;
  justify-content: center;
}
.profile-title-modern {
  font-size: 1.28rem;
  font-weight: 900;
  color: #212942;
}
.profile-role-modern {
  color: #2563eb;
  font-size: 1.09rem;
  font-weight: 700;
  opacity: .89;
  margin-top: .08rem;
}

.profile-info-modern {
  background: #fff;
  border-radius: 12px;
  padding: 1.1rem 1.25rem 1.1rem 1.25rem;
  margin-bottom: .55rem;
  margin-top: 0.2rem;
  display: flex;
  flex-direction: column;
  gap: 1.03rem;
  font-size: 1.05rem;
  box-shadow: 0 2.5px 14px #e3eaff19;
}
.info-row {
  display: flex;
  align-items: flex-start;
  gap: 0.8rem;
  line-height: 1.45;
}
.info-label {
  color: #7082a4;
  font-weight: 700;
  font-size: 1.04rem;
  min-width: 85px;
  display: flex;
  align-items: center;
  letter-spacing: .01em;
}
.info-value {
  color: #1a2544;
  font-weight: 800;
  font-size: 1.08rem;
  letter-spacing: .02em;
}

.profile-divider {
  height: 1px;
  background: linear-gradient(90deg, #d4e0f6 10%, #fff 80%);
  border: none;
  margin: .55rem 0 .35rem 0;
}

.profile-change-title-modern {
  font-size: 1.05rem;
  font-weight: 800;
  color: #293052;
  margin-top: .5rem;
  margin-bottom: .16rem;
}

.profile-change-fields-modern {
  display: flex;
  flex-direction: column;
  margin-top: .3rem;
  margin-bottom: .35rem;
  gap: .69rem;
}
.input-modern {
  font-size: 1.08rem;
  border-radius: 8px;
  border: 1.1px solid #dde8f7;
  padding: 0.48rem 1.08rem;
  background: #f5f7fb;
  transition: border .13s, box-shadow .13s;
}
.input-modern:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 2px #2563eb22;
  outline: none;
}

.profile-actions-modern {
  display: flex;
  gap: .7rem;
  justify-content: flex-end;
  margin-top: .08rem;
}

.change-btn-modern {
  background: linear-gradient(93deg, #2563eb 80%, #60a5fa 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 0.52rem 1.7rem;
  font-size: 1rem;
  font-weight: 900;
  cursor: pointer;
  transition: background 0.13s, transform 0.11s;
  letter-spacing: .01em;
  box-shadow: 0 1.5px 8px #4f8cff11;
}
.change-btn-modern:hover {
  background: linear-gradient(93deg, #1749d4 85%, #2563eb 100%);
  transform: scale(1.035);
}

.close-btn-modern {
  background: #f4f5fa;
  color: #1a2544;
  border: none;
  border-radius: 8px;
  padding: 0.52rem 1.25rem;
  font-size: 1rem;
  font-weight: 700;
  transition: background .14s, color .12s;
}
.close-btn-modern:hover {
  background: #e6edff;
  color: #1c2241;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.18s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

.error-msg-modern, .success-msg-modern {
  font-size: 1.06rem;
  text-align: center;
  font-weight: 700;
  margin-top: 0.18rem;
  letter-spacing: .01em;
}
.error-msg-modern { color: #e11d48; }
.success-msg-modern { color: #22c55e; }

@media (max-width: 520px) {
  .profile-modal-modern {
    width: 99vw;
    padding: 1.2rem .2rem 1rem .2rem;
    min-height: unset;
    border-radius: 11px;
  }
  .profile-info-modern {
    padding: 0.7rem 0.7rem 0.7rem 0.7rem;
  }
  .profile-header-modern {
    gap: .6rem;
  }
}
.profile-info-modern-rich {
  background: rgba(248, 252, 255, 0.96);
  border-radius: 16px;
  border: 1.6px solid #e3eaff;
  box-shadow: 0 4px 24px #dde8f74d;
  padding: 1.12rem 1.25rem 1.12rem 1.25rem;
  margin-bottom: 1.15rem;
  margin-top: 0.28rem;
  display: flex;
  flex-direction: column;
  gap: 1.09rem;
  font-size: 1.11rem;
  position: relative;
  backdrop-filter: blur(3.5px);
  transition: box-shadow 0.22s;
}

.profile-info-modern-rich:hover {
  box-shadow: 0 8px 32px #b7c5ec2e, 0 2px 10px #2563eb11;
  border-color: #b7c5ec;
}

.info-row-rich {
  display: flex;
  align-items: center;
  gap: 0.74rem;
  line-height: 1.42;
}

.info-icon-circle {
  background: linear-gradient(120deg, #e3eaff 55%, #f3f8ff 100%);
  border-radius: 50%;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 4px;
  box-shadow: 0 1.5px 6px #2563eb0f;
}

.info-label-rich {
  color: #6c81a8;
  font-size: 1.03rem;
  font-weight: 700;
  min-width: 54px;
  margin-right: 0.2em;
  letter-spacing: 0.01em;
}

.info-value-rich {
  color: #203050;
  font-weight: 900;
  font-size: 1.14rem;
  margin-left: 0.2em;
  letter-spacing: .01em;
}
.profile-change-fields-modern {
  display: flex;
  flex-direction: column;
  margin-top: .22rem;
  margin-bottom: .15rem;
  gap: .32rem;               /* стало компактнее */
}

.input-modern {
  font-size: 1.08rem;
  border-radius: 10px;
  border: 1.2px solid #e3eaff;
  padding: 0.64rem 1.15rem;  /* стало чуть компактнее */
  background: linear-gradient(120deg, #f9fbff 70%, #ecf3fa 100%);
  box-shadow: 0 1.5px 8px #d9eaff10;
  margin: 0;
  transition: border .15s, box-shadow .18s, background .18s;
  outline: none;
}

.input-modern:focus {
  border: 1.4px solid #2563eb;
  background: #f0f6ff;
  box-shadow: 0 2px 12px #2563eb18;
}

.input-modern::placeholder {
  color: #92a3c5;
  font-weight: 600;
  opacity: 1;
  letter-spacing: 0.01em;
  font-size: 1.03rem;
}
.accent-avatar {
  background: linear-gradient(120deg, #a5b6fa 75%, #2563eb 100%);
  box-shadow: 0 2.5px 16px #2563eb26, 0 2.5px 10px #2563eb18;
}

.user-details-compact {
  gap: .58rem;
  padding: 0.7rem 1.05rem 0.85rem 1.05rem;
  border-width: 1.8px;
}

.main-name {
  font-size: 1.35rem;
  font-weight: 900;
  letter-spacing: .01em;
}

.input-shadow {
  box-shadow: 0 2px 12px #b7c5ec0d;
  border-radius: 10px;
}

.input-modern {
  font-size: 1.03rem;
  padding: 0.53rem 1.05rem;
}

.profile-change-fields-modern {
  gap: .38rem;
}

.main-btn-strong {
  background: linear-gradient(95deg, #2563eb 70%, #4785ff 100%);
  box-shadow: 0 3px 16px #2563eb18;
  font-weight: 900;
  font-size: 1.04rem;
}

.main-btn-strong:hover {
  background: linear-gradient(93deg, #1749d4 85%, #2563eb 100%);
  transform: scale(1.038);
}

.main-btn-ghost {
  border: 1.6px solid #d7e4fa;
  background: #f4f8fd;
  color: #28385e;
  font-weight: 700;
  transition: background .15s, color .13s;
}

.main-btn-ghost:hover {
  background: #e9f0ff;
  color: #133179;
  border-color: #b7c5ec;
}

/* Слегка уменьшить расстояния */
.profile-modal-modern {
  gap: 1.1rem;
  padding-top: 2.0rem;
  padding-bottom: 1.4rem;
}

@media (max-width: 520px) {
  .profile-modal-modern { padding: 0.7rem 0.07rem 0.8rem 0.07rem; }
}
.profile-card-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

/* Общий стиль карточки для обоих блоков */
.profile-info-modern-rich {
  background: #fff;
  border-radius: 14px;
  border: 1.5px solid #e3eaff;
  box-shadow: 0 4px 24px #dde8f74d;
  padding: 1.08rem 1.3rem 1.08rem 1.3rem;
  display: flex;
  flex-direction: column;
  gap: 1.08rem;
}

/* Дополнительная компактность для информации */
.profile-info-compact {
  gap: .63rem;
  padding-bottom: .73rem;
}

.profile-password-compact {
  margin-top: 0; /* сбросить лишний отступ */
  gap: 0.8rem;
  padding-top: 1rem;
}

/* Группа кнопок ниже! */
.profile-actions-modern.buttons-bottom {
  margin-top: 1.2rem;
  justify-content: flex-end;
}

/* Уменьшаем контраст разделителя */
.profile-divider {
  display: none;
}
.profile-change-fields-modern.profile-fields-spaced {
  display: flex;
  flex-direction: column;
  gap: 1.15rem;   /* Увеличенный отступ между инпутами */
}

.profile-actions-modern.profile-actions-outside {
  margin-top: 2.2rem;
  justify-content: flex-end;
  gap: 1.1rem;
  display: flex;
  /* можно добавить паддинг вниз, если нужно */
}

/* Для мобильных — чтобы не улетало за пределы */
@media (max-width: 520px) {
  .profile-actions-modern.profile-actions-outside {
    margin-top: 1.4rem;
  }
}

</style>

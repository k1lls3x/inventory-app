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
          <div class="dashboard-risk-grid">
            <div class="table-section dashboard-risk-card">
              <div class="table-header dashboard-risk-header">
                <div>
                  <p class="title">Контроль дефицита</p>
                  <p class="dashboard-subtitle">Товары ниже минимального остатка и бюджет на пополнение</p>
                </div>
                <select v-model.number="selectedRiskWarehouseId" class="input dashboard-risk-filter">
                  <option :value="0">Все склады</option>
                  <option v-for="wh in warehouses" :key="wh.warehouse_id" :value="wh.warehouse_id">
                    {{ wh.name }}
                  </option>
                </select>
              </div>

              <div class="risk-summary-cards">
                <div class="risk-summary-card warning">
                  <span class="risk-summary-label">Рисковых SKU</span>
                  <strong>{{ lowStockItems.length }}</strong>
                </div>
                <div class="risk-summary-card danger">
                  <span class="risk-summary-label">Критичных позиций</span>
                  <strong>{{ criticalLowStockCount }}</strong>
                </div>
                <div class="risk-summary-card neutral">
                  <span class="risk-summary-label">К дозакупке</span>
                  <strong>{{ recommendedRestockUnits }}</strong>
                </div>
                <div class="risk-summary-card accent">
                  <span class="risk-summary-label">Оценка бюджета</span>
                  <strong>{{ formatMoney(recommendedRestockCost) }}</strong>
                </div>
              </div>

              <table>
                <thead>
                  <tr>
                    <th>Товар</th>
                    <th>Склад</th>
                    <th>Остаток</th>
                    <th>Мин. остаток</th>
                    <th>Рекомендация</th>
                    <th>Статус</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in lowStockItems.slice(0, 8)" :key="`${item.item_id}-${item.warehouse_id}`">
                    <td>
                      <div>{{ item.name }}</div>
                      <div class="dashboard-row-meta">{{ item.sku }}</div>
                    </td>
                    <td>{{ item.warehouse }}</td>
                    <td>{{ item.current_stock }}</td>
                    <td>{{ item.reorder_level }}</td>
                    <td>{{ item.recommended_order }}</td>
                    <td>
                      <span class="stock-health-badge" :class="item.severity">
                        {{ item.status_label }}
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>

              <div v-if="lowStockItems.length === 0" class="empty-message">
                Все позиции находятся выше минимального остатка
              </div>
            </div>

            <div class="chart-card dashboard-risk-sidecard">
              <p class="title">Склады под нагрузкой</p>
              <div class="warehouse-risk-list">
                <div v-for="item in warehouseRiskSummary" :key="item.warehouse_id" class="warehouse-risk-row">
                  <div>
                    <div class="warehouse-risk-name">{{ item.name }}</div>
                    <div class="dashboard-row-meta">{{ item.low_count }} SKU требуют внимания</div>
                  </div>
                  <div class="warehouse-risk-values">
                    <strong>{{ item.critical_count }}</strong>
                    <span>критично</span>
                  </div>
                </div>
              </div>
              <div v-if="warehouseRiskSummary.length === 0" class="empty-message">
                По текущим данным дефицита нет
              </div>
            </div>
          </div>
        </section>
 <!-- Пользователи (видно только админу) -->
 <section v-if="currentTab === 'Пользователи' && user?.role === 'admin'">
  <div class="filters-bar">
    <div class="filter-group">
      <label>🔍 Поиск</label>
      <input
        type="text"
        class="input"
        v-model="userSearch"
        placeholder="Имя пользователя, ФИО или роль"
      />
    </div>
    <div class="filter-group button-group">
      <label>&nbsp;</label>
      <button class="add-button" @click="openAddUserModal">➕ Добавить пользователя</button>
    </div>
  </div>

  <div class="cards">
    <div class="card animate-card">
      <p class="title">Всего пользователей</p>
      <p class="value">{{ users.length }}</p>
    </div>
  </div>

  <div class="table-section animate-table">
    <div class="table-header">
      <p class="title">Пользователи</p>
      <button class="export-button" @click="exportUsersToExcel">📤 Экспорт в Excel</button>
    </div>
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Логин</th>
          <th>ФИО</th>
          <th>Роль</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="u in filteredUsers" :key="u.user_id">
          <td>{{ u.user_id }}</td>
          <td>{{ u.username }}</td>
          <td>{{ u.full_name }}</td>
          <td>{{ roleName(u.role) }}</td>
          <td>
            <div class="action-buttons">
              <button class="action-btn edit" @click="openEditUserModal(u)">✏️</button>
              <button class="action-btn delete" @click="deleteUser(u)">🗑️</button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="filteredUsers.length === 0" class="empty-message">
      Нет пользователей по фильтру
    </div>
  </div>

  <!-- Модалка добавления пользователя -->
  <div v-if="showAddUserModal" class="modal-overlay" @click.self="closeAddUserModal">
    <div class="modal">
      <h3>Добавить пользователя</h3>
      <div class="form-group"><label>Логин</label><input v-model="newUser.username" /></div>
      <div class="form-group"><label>ФИО</label><input v-model="newUser.full_name" /></div>
      <div class="form-group"><label>Пароль</label><input type="password" v-model="newUser.password" /></div>
      <div class="form-group"><label>Роль</label>
        <select v-model="newUser.role">
          <option disabled value="">Выберите роль</option>
          <option value="admin">Администратор</option>
          <option value="manager">Менеджер</option>
          <option value="worker">Сотрудник</option>
        </select>
      </div>
      <div class="modal-actions">
        <button @click="confirmAddUser">💾 Сохранить</button>
        <button @click="closeAddUserModal">❌ Отмена</button>
      </div>
    </div>
  </div>

  <!-- Модалка редактирования пользователя -->
  <div v-if="showEditUserModal" class="modal-overlay" @click.self="closeEditUserModal">
    <div class="modal">
      <h3>Редактировать пользователя</h3>
      <div class="form-group"><label>Логин</label><input v-model="userToEdit.username" /></div>
      <div class="form-group"><label>ФИО</label><input v-model="userToEdit.full_name" /></div>
      <div class="form-group"><label>Новый пароль</label><input type="password" v-model="userToEdit.newPassword" placeholder="Оставьте пустым для без изменений" /></div>
      <div class="form-group"><label>Роль</label>
        <select v-model="userToEdit.role">
          <option value="admin">Администратор</option>
          <option value="manager">Менеджер</option>
          <option value="worker">Сотрудник</option>
        </select>
      </div>
      <div class="modal-actions">
        <button @click="confirmEditUser">💾 Сохранить</button>
        <button @click="closeEditUserModal">❌ Отмена</button>
      </div>
    </div>
  </div>
</section>

        <!-- Остатки -->
        <section v-if="currentTab === 'Остатки' && ['admin', 'manager', 'worker'].includes(user?.role)">
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
        <section v-if="currentTab === 'Поставки' && ['admin', 'manager', 'worker'].includes(user?.role)">          <div class="filters-bar">
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
                  <option value="-1">➕ Новый товар...</option>
                  <option v-for="item in items" :key="item.item_id" :value="item.item_id">
                    {{ item.name }} ({{ item.sku }})
                  </option>
                </select>
                <div v-if="newInbound.item_id === -1" class="new-item-fields">
                  <input class="input" placeholder="SKU" v-model="newInboundItem.sku" />
                  <input class="input" placeholder="Наименование" v-model="newInboundItem.name" />
                  <input class="input" placeholder="Описание" v-model="newInboundItem.description" />
                  <input class="input" placeholder="Ед. изм." v-model="newInboundItem.uom" />
                  <input class="input" type="number" placeholder="Цена" v-model.number="newInboundItem.price" />
                  <input class="input" type="number" placeholder="Себестоимость" v-model.number="newInboundItem.cost" />
                </div>
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

        </section>

        <!-- Товары -->
        <section v-if="currentTab === 'Товары' && ['admin', 'manager'].includes(user?.role)">
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



        </section>
        <section v-if="currentTab === 'Склады'">
  <!-- Карточки -->
  <div class="cards">
    <div class="card animate-card">
      <p class="title">Всего складов</p>
      <p class="value">{{ warehouses.length }}</p>
    </div>
    <!-- Можно добавить другие карточки по необходимости -->
  </div>

  <!-- Фильтры -->
  <div class="filters-bar">
    <div class="filter-group">
      <label>🔍 Поиск</label>
      <input
        type="text"
        class="input"
        v-model="warehouseSearch"
        placeholder="Название или локация"
      />
    </div>
    <div class="filter-group button-group">
      <label>&nbsp;</label>
      <button class="add-button" @click="openAddModal">➕ Добавить склад</button>
    </div>
  </div>

  <!-- Таблица -->
  <div class="table-section animate-table">
    <div class="table-header">
      <p class="title">Склады</p>
      <!-- <button class="export-button">📤 Экспорт</button> -->
    </div>
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Название</th>
          <th>Локация</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="w in filteredWarehouses" :key="w.warehouse_id">
          <td>{{ w.warehouse_id }}</td>
          <td>{{ w.name }}</td>
          <td>{{ w.location }}</td>
          <td>
            <div class="action-buttons">
              <button class="action-btn edit" @click="editWarehouse(w)">✏️</button>
              <!-- Тут не делаем кнопку “Удалить”, т.к. удаление складов нежелательно -->
            </div>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="filteredWarehouses.length === 0" class="empty-message">
      Нет складов по фильтру
    </div>
  </div>

  <!-- Модалка добавления склада -->
  <div v-if="showAddModal" class="modal-overlay" @click.self="closeAddModal">
    <div class="modal">
      <h3>Добавить склад</h3>
      <div class="form-group">
        <label>Название склада</label>
        <input v-model="newWarehouse.name" />
      </div>
      <div class="form-group">
        <label>Локация</label>
        <input v-model="newWarehouse.location" />
      </div>
      <div class="modal-actions">
        <button @click="addWarehouse">💾 Сохранить</button>
        <button @click="closeAddModal">❌ Отмена</button>
      </div>
    </div>
  </div>

  <!-- Модалка редактирования склада -->
  <div v-if="showEditModal" class="modal-overlay" @click.self="closeEditModal">
    <div class="modal">
      <h3>Редактировать склад</h3>
      <div class="form-group">
        <label>Название склада</label>
        <input v-model="editWarehouseData.name" />
      </div>
      <div class="form-group">
        <label>Локация</label>
        <input v-model="editWarehouseData.location" />
      </div>
      <div class="modal-actions">
        <button @click="updateWarehouse">💾 Сохранить</button>
        <button @click="closeEditModal">❌ Отмена</button>
      </div>
    </div>
  </div>
</section>
<section v-if="currentTab === 'Движения'">
  <div class="filters-bar">
    <div class="filter-group">
      <label>Тип движения</label>
      <select v-model="moveType" class="input">
        <option value="">Все</option>
        <option value="inbound">Поступление</option>
        <option value="outbound">Отгрузка</option>
      </select>
    </div>
    <div class="filter-group">
      <label>📦 Склад</label>
      <select v-model.number="moveWarehouseId" class="input">
        <option value="0">Все склады</option>
        <option v-for="wh in warehouses" :key="wh.warehouse_id" :value="wh.warehouse_id">{{ wh.name }}</option>
      </select>
    </div>
    <div class="filter-group">
      <label>🔍 Поиск товара</label>
      <input v-model="moveItemSearch" placeholder="Название или SKU" class="input" />
    </div>
  </div>

  <div class="cards">
    <div class="card animate-card">
      <p class="title">Всего поступлений</p>
      <p class="value">{{ inboundCount }}</p>
    </div>
    <div class="card animate-card">
      <p class="title">Всего отгрузок</p>
      <p class="value">{{ outboundCount }}</p>
    </div>
    <div class="card animate-card">
      <p class="title">Общий оборот</p>
      <p class="value">{{ totalMoved }}</p>
    </div>
  </div>

  <div class="charts-table-wrap">
    <div class="chart-card animate-chart">
      <LineChart :data="moveChartData" />
    </div>
  </div>

  <div class="table-section animate-table">
    <div class="table-header">
      <p class="title">Движения товаров</p>
    </div>
    <table>
      <thead>
        <tr>
          <th>Дата</th>
          <th>Тип</th>
          <th>Товар</th>
          <th>Склад</th>
          <th>Кол-во</th>
          <th v-if="moveType === 'inbound' || !moveType">Поставщик</th>
          <th v-if="moveType === 'outbound' || !moveType">Описание</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="m in filteredMovements" :key="m.movement_id">
  <td>{{ formatDate(m.date) }}</td>
  <td>{{ movementTypeName(m.type) }}</td>
  <td>{{ m.item_name }}</td>
  <td>{{ m.warehouse_name }}</td>
  <td :class="{ 'positive': m.type==='inbound','negative': m.type==='outbound' }">
    {{ m.quantity }}
  </td>
  <!-- всегда рендерим две ячейки в одном порядке -->
  <td>
    {{ m.type === 'inbound' ? (m.supplier_name || '—') : '—' }}
  </td>
  <td>
    {{ m.type === 'outbound' ? (m.destination     || '—') : '—' }}
  </td>
</tr>
      </tbody>
    </table>
    <div v-if="filteredMovements.length === 0" class="empty-message">Нет движений по фильтру</div>
  </div>
</section>

<!-- Поставщики -->
<section v-if="currentTab === 'Поставщики' && ['admin', 'manager'].includes(user?.role)">


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
    <div class="form-group"><label>Телефон</label> <input
          v-model="newSupplier.phone"
          @input="maskPhone($event, supplierToEdit)"
          maxlength="18"
          placeholder="+7 (___)-___-__-__"
          type="tel"
        /></div>
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
    <div class="form-group"><label>Телефон</label><input
  v-model="newSupplier.phone"
  @input="maskPhone($event, supplierToEdit)"
  maxlength="18"
  placeholder="+7 (___)-___-__-__"
  type="tel"
/></div>
    <div class="form-group"><label>Email</label><input v-model="supplierToEdit.email" /></div>
    <div class="modal-actions">
      <button @click="confirmEditSupplier">💾 Сохранить</button>
      <button @click="showEditSupplierModal = false">❌ Отмена</button>
    </div>
  </div>
</div>
</section>
<section v-if="currentTab === 'Отгрузки'">
  <div class="filters-bar">
    <div class="filter-group">
      <label>📅 Дата</label>
      <input type="date" v-model="outboundDateFilter" class="input" />
    </div>
    <div class="filter-group">
      <label>🔍 Поиск</label>
      <input v-model="outboundSearch" class="input" placeholder="Товар, SKU, склад, получатель" />
    </div>
    <div class="filter-group button-group">
      <label>&nbsp;</label>
      <button class="add-button" @click="openAddOutboundModal">➕ Добавить отгрузку</button>
    </div>
  </div>

  <div class="cards">
    <div class="card animate-card">
      <p class="title">Всего отгрузок</p>
      <p class="value">{{ outboundList.length }}</p>
    </div>
    <div class="card animate-card">
      <p class="title">Суммарно отгружено</p>
      <p class="value">{{ totalOutboundQuantity }}</p>
    </div>
  </div>

  <div class="charts-table-wrap">
    <div class="chart-card animate-chart">
      <LineChart :data="outboundChartData" />
    </div>
  </div>

  <div class="table-section animate-table">
    <div class="table-header">
      <p class="title">Отгрузки</p>
    </div>
    <table>
      <thead>
        <tr>
          <th>Дата</th>
          <th>Товар</th>
          <th>SKU</th>
          <th>Склад</th>
          <th>Описание</th>
          <th>Кол-во</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="o in filteredOutboundList" :key="o.outbound_id">
          <td>{{ formatDate(o.date) }}</td>
          <td>{{ o.name }}</td>
          <td>{{ o.sku }}</td>
          <td>{{ o.warehouse }}</td>
          <td>{{ o.destination }}</td>
          <td class="negative">{{ o.quantity }}</td>
          <td>
            <div class="action-buttons">
              <button class="action-btn edit" @click="openEditOutboundModal(o)">✏️</button>
              <button class="action-btn delete" @click="deleteOutbound(o)">🗑️</button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="filteredOutboundList.length === 0" class="empty-message">
      Нет отгрузок по фильтру
    </div>
  </div>

  <!-- Модалка добавления -->
  <div v-if="showAddOutboundModal" class="modal-overlay" @click.self="closeAddOutboundModal">
    <div class="modal">
      <h3>Добавить отгрузку</h3>
      <div class="form-group"><label>Товар</label>
        <select v-model.number="newOutbound.item_id">
          <option disabled value="0">Выберите товар</option>
          <option v-for="item in items" :key="item.item_id" :value="item.item_id">{{ item.name }} ({{ item.sku }})</option>
        </select>
      </div>
      <div class="form-group"><label>Склад</label>
        <select v-model.number="newOutbound.warehouse_id">
          <option disabled value="0">Выберите склад</option>
          <option v-for="w in warehouses" :key="w.warehouse_id" :value="w.warehouse_id">{{ w.name }}</option>
        </select>
      </div>
      <div class="form-group"><label>Описание</label>
        <input v-model="newOutbound.destination" />
      </div>
      <div class="form-group"><label>Количество</label>
        <input type="number" min="1" v-model.number="newOutbound.quantity" />
      </div>
      <div class="form-group"><label>Дата отгрузки</label>
        <input type="date" v-model="newOutbound.shipped_at" />
      </div>
      <div class="modal-actions">
        <button @click="confirmAddOutbound">💾 Сохранить</button>
        <button @click="closeAddOutboundModal">❌ Отмена</button>
      </div>
    </div>
  </div>

  <!-- Модалка редактирования -->
  <div v-if="showEditOutboundModal" class="modal-overlay" @click.self="closeEditOutboundModal">
    <div class="modal">
      <h3>Редактировать отгрузку</h3>
      <div class="form-group"><label>Товар</label>
        <select v-model.number="outboundToEdit.item_id">
          <option disabled value="0">Выберите товар</option>
          <option v-for="item in items" :key="item.item_id" :value="item.item_id">{{ item.name }} ({{ item.sku }})</option>
        </select>
      </div>
      <div class="form-group"><label>Склад</label>
        <select v-model.number="outboundToEdit.warehouse_id">
          <option disabled value="0">Выберите склад</option>
          <option v-for="w in warehouses" :key="w.warehouse_id" :value="w.warehouse_id">{{ w.name }}</option>
        </select>
      </div>
      <div class="form-group"><label>Описание</label>
        <input v-model="outboundToEdit.destination" />
      </div>
      <div class="form-group"><label>Количество</label>
        <input type="number" min="1" v-model.number="outboundToEdit.quantity" />
      </div>
      <div class="form-group"><label>Дата отгрузки</label>
        <input type="date" v-model="outboundToEdit.shipped_at" />
      </div>
      <div class="modal-actions">
        <button @click="confirmEditOutbound">💾 Сохранить</button>
        <button @click="closeEditOutboundModal">❌ Отмена</button>
      </div>
    </div>
  </div>
</section>
        <!-- Другое (заглушка) -->
        <section v-if="!tabs.includes(currentTab)">
    <p>Раздел "{{ currentTab }}" в разработке или доступ ограничен...</p>
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
import { ExportUsersToExcel, ExportSuppliersToExcel, ExportDeliveriesToExcel, ExportItemsToExcel } from '../wailsjs/go/app/App'
import LoginForm from './components/LoginForm.vue'
import {
  GetOutboundDetails, AddOutbound, EditOutbound, RemoveOutbound
} from '../wailsjs/go/app/App'
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
  AddInboundTx,
  GetSuppliers,
  DeleteInbound,
  EditInbound,
  ChangePassword,
  GetItems,
  AddItem,
  UpdateItem,
  RemoveItem,
  EditSupplier,
  AddSupplier,
  RemoveSupplier,
  GetUsers,
  RegisterUser,
  RemoveUser,
  ChangeUserData,
  AddWarehouse,EditWarehouse,
  GetAllMovementsThisMonth

} from '../wailsjs/go/app/App'

const loggedIn = ref(localStorage.getItem('loggedIn') === 'true')
const emit = defineEmits(['login-success'])

async function onLoginSuccess(userData) {
  user.value = userData     
  loggedIn.value = true
  localStorage.setItem('loggedIn', 'true')
  await loadInitialData()
}

async function loadInitialData() {
  items.value = await GetItems() || []
  weeklyStockData.value = await GetWeeklyStockTrend() || []
  const dash = await GetDashboard()
  totalStock.value = dash.total_stock
  itemCount.value = dash.item_count
  monthlyOrders.value = dash.monthly_orders
  newItems.value = dash.new_items
  suppliers.value = await GetSuppliers() || []
  const stockData = normalizeStockRows(await GetStockDetails())
  stockList.value = stockData
  allStockList.value = stockData
  await reloadMovements()
  await reloadOutbound()
  try {
    deliveriesList.value = await GetInboundDetails() || []
  } catch (err) {
    console.error('Ошибка загрузки поставок:', err)
  }
  topItems.value = await GetTopItems() || []
  await loadWarehouses()
  turnoverData.value = await GetTurnoverByWarehouse() || []
}
function logout() {
  localStorage.removeItem('loggedIn')
  loggedIn.value = false
}

// Пользователи
const users = ref([])
const userSearch = ref('')
const showAddUserModal = ref(false)
const showEditUserModal = ref(false)
const newUser = ref({ username: '', full_name: '', password: '', role: '' })
const userToEdit = ref({})

const showAddItemModal = ref(false)
const showEditItemModal = ref(false)
const itemToEdit = ref(null)
// Состояния
const outboundList = ref([])
const showAddOutboundModal = ref(false)
const showEditOutboundModal = ref(false)
const outboundToEdit = ref({})
const outboundSearch = ref('')
const outboundDateFilter = ref('')
const newOutbound = ref({
  item_id: 0, warehouse_id: 0, destination: '', quantity: 1, shipped_at: ''
})

const filteredUsers = computed(() =>
  users.value.filter(u =>
    (u.username || '').toLowerCase().includes(userSearch.value.toLowerCase()) ||
    (u.full_name || '').toLowerCase().includes(userSearch.value.toLowerCase()) ||
    (roleName(u.role) || '').toLowerCase().includes(userSearch.value.toLowerCase())
  )
)

function openAddUserModal() { showAddUserModal.value = true }
function closeAddUserModal() {
  showAddUserModal.value = false
  newUser.value = { username: '', full_name: '', password: '', role: '' }
}

async function confirmAddUser() {
  if (!newUser.value.username || !newUser.value.full_name || !newUser.value.password || !newUser.value.role) {
    alert('Заполните все поля')
    return
  }
  try {
    await RegisterUser(
  newUser.value.username,   // первый аргумент: username
  newUser.value.password,   // второй аргумент: password
  newUser.value.full_name,  // третий аргумент: fullName
  newUser.value.role        // четвертый аргумент: role
)
    users.value = await GetUsers() || []
    closeAddUserModal()
  } catch (e) {
    alert('Ошибка при добавлении пользователя: ' + (e?.message || ''))
  }
}


function openEditUserModal(u) {
  userToEdit.value = { ...u, newPassword: '' }
  showEditUserModal.value = true
}

function closeEditUserModal() { showEditUserModal.value = false }
async function confirmEditUser() {
  if (!userToEdit.value.username || !userToEdit.value.full_name || !userToEdit.value.role) {
    alert('Заполните все поля')
    return
  }
  try {
    const payload = {
      user_id: userToEdit.value.user_id,
      username: userToEdit.value.username,
      full_name: userToEdit.value.full_name,
      role: userToEdit.value.role,
      // Передавай новый пароль только если он не пустой
      password: userToEdit.value.newPassword || undefined
    }
    await ChangeUserData(payload)
    users.value = await GetUsers() || []
    closeEditUserModal()
  } catch (e) {
    alert('Ошибка при обновлении пользователя: ' + (e?.message || ''))
  }
}
const movements = ref([])
const moveType = ref("")
const moveWarehouseId = ref(0)
const moveItemSearch = ref("")

// График - динамика по дням (поступления и отгрузки)
const moveChartData = computed(() => {
  // Группируем движения по дате, считаем поступления/отгрузки отдельно
  const byDate = {}
  for (const m of movements.value) {
    const date = (new Date(m.date)).toISOString().slice(0, 10)
    if (!byDate[date]) byDate[date] = { in: 0, out: 0 }
    if (m.type === 'inbound') byDate[date].in += m.quantity
    if (m.type === 'outbound') byDate[date].out += m.quantity
  }
  const dates = Object.keys(byDate).sort()
  return {
    labels: dates.map(d => d.split('-').reverse().join('.')),
    datasets: [
      {
        label: "Поступления",
        data: dates.map(d => byDate[d].in),
        borderColor: "#22c55e",
        backgroundColor: "rgba(34,197,94,0.1)",
        tension: 0.3
      },
      {
        label: "Отгрузки",
        data: dates.map(d => byDate[d].out),
        borderColor: "#ef4444",
        backgroundColor: "rgba(239,68,68,0.1)",
        tension: 0.3
      }
    ]
  }
})

const filteredMovements = computed(() =>
  movements.value.filter(m => {
    const wid = Number(moveWarehouseId.value)  // ← здесь
    return (!moveType.value || m.type === moveType.value)
        && (wid === 0 || m.warehouse_id === wid)
        && (
            m.item_name.toLowerCase().includes(moveItemSearch.value.toLowerCase()) ||
            m.item_id.toString().includes(moveItemSearch.value.toLowerCase())
        )
  })
)


const inboundCount  = computed(() =>
  filteredMovements.value.filter(m => m.type === 'inbound').length
)

const outboundCount = computed(() =>
  filteredMovements.value.filter(m => m.type === 'outbound').length)
const totalMoved = computed(() => movements.value.reduce((acc, m) => acc + m.quantity, 0))

function movementTypeName(t) {
  if (t === 'inbound') return 'Поступление'
  if (t === 'outbound') return 'Отгрузка'
  return t
}


async function reloadMovements() {
  // Можно сюда добавить аргументы фильтров, если делаешь серверный фильтр
  movements.value = await GetAllMovementsThisMonth() || []
}

async function deleteUser(u) {
  if (!confirm(`Удалить пользователя "${u.username}"?`)) return
  try {
    await RemoveUser(u.user_id)
    users.value = await GetUsers() || []
  } catch (e) {
    alert('Ошибка при удалении: ' + (e?.message || ''))
  }
}

function exportUsersToExcel() {
  ExportUsersToExcel().then(base64data => {
    const binary = atob(base64data)
    const len = binary.length
    const bytes = new Uint8Array(len)
    for (let i = 0; i < len; i++) {
      bytes[i] = binary.charCodeAt(i)
    }
    const blob = new Blob([bytes], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = 'users.xlsx'
    link.click()
    setTimeout(() => URL.revokeObjectURL(link.href), 1000)
  }).catch(err => {
    alert('Ошибка экспорта: ' + err)
  })
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
      'Отгрузки',
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
      'Отгрузки',
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
    'Отгрузки',
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
const allStockList = ref([])
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
const newInboundItem = ref({
  sku: '',
  name: '',
  description: '',
  uom: '',
  price: 0,
  cost: 0,
})
const suppliers = ref([])
const selectedRiskWarehouseId = ref(0)
 // если вкладки реализованы через состояние, замени на свой способ
const warehouseSearch = ref('')
const newWarehouse = ref({ name: '', location: '' })
const editWarehouseData = ref({ warehouse_id: null, name: '', location: '' })
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

async function loadWarehouses() {
  warehouses.value = await GetWarehouses() || []
}


function openAddSupplierModal() { showAddSupplierModal.value = true }
async function confirmAddSupplier() {
  // Валидация при необходимости
  try {
    await window.go.app.App.AddSupplier(newSupplier.value)
    showAddSupplierModal.value = false
    newSupplier.value = { name: '', inn: '', contact_person: '', phone: '', email: '' }
    // После добавления — обновить список с бэка
    suppliers.value = await GetSuppliers() || []
  } catch (e) {
    alert('Ошибка при добавлении: ' + (e?.message || ''))
  }
}

async function confirmEditSupplier() {
  try {
    await window.go.app.App.EditSupplier(supplierToEdit.value)
    showEditSupplierModal.value = false
    suppliers.value = await GetSuppliers() || []
  } catch (e) {
    alert('Ошибка при обновлении: ' + (e?.message || ''))
  }
}

function openEditSupplierModal(s) { supplierToEdit.value = { ...s }; showEditSupplierModal.value = true }

async function deleteSupplier(s) {
  if (!confirm(`Удалить поставщика "${s.name}"?`)) return
  try {
    // Удаляем по уникальному supplier_id
    await window.go.app.App.RemoveSupplier(s.supplier_id)
    // После удаления — обновить список с бэка
    suppliers.value = await GetSuppliers() || []
  } catch (e) {
    alert('Ошибка при удалении: ' + (e?.message || ''))
  }
}

function maskPhone(event, obj) {
  let v = event.target.value.replace(/\D/g, '');
  if (v.startsWith('8')) v = '7' + v.slice(1); // Преобразуем 8 -> 7
  if (!v.startsWith('7')) v = '7' + v;
  v = v.slice(0, 11);

  let res = '+7';
  if (v.length > 1) res += ' (' + v.slice(1, 4);
  if (v.length >= 4) res += ')';
  if (v.length >= 4) res += '-' + v.slice(4, 7);
  if (v.length >= 7) res += '-' + v.slice(7, 9);
  if (v.length >= 9) res += '-' + v.slice(9, 11);
  obj.phone = res;
}


function exportSuppliersToExcel() {
  ExportSuppliersToExcel().then(base64data => {
    const binary = atob(base64data)
    const len = binary.length
    const bytes = new Uint8Array(len)
    for (let i = 0; i < len; i++) {
      bytes[i] = binary.charCodeAt(i)
    }
    const blob = new Blob([bytes], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = 'suppliers.xlsx'
    link.click()
    setTimeout(() => URL.revokeObjectURL(link.href), 1000)
  }).catch(err => {
    alert('Ошибка экспорта: ' + err)
  })
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
  newInbound.value = { item_id: 0, supplier_id: 0, warehouse_id: 0, quantity: 1, received_at: '' }
  Object.assign(newInboundItem.value, { sku: '', name: '', description: '', uom: '', price: 0, cost: 0 })
}

function normalizeStockRows(rows) {
  return (rows || []).map(s => ({
    id: s.stock_id,
    stock_id: s.stock_id,
    item_id: s.item_id,
    warehouse_id: s.warehouse_id,
    name: s.name,
    sku: s.sku,
    warehouse: s.warehouse,
    quantity: s.quantity
  }))
}

async function refreshAllStockDetails() {
  const rows = normalizeStockRows(await GetStockDetails())
  allStockList.value = rows
  if (Number(selectedWarehouseId.value) === 0) {
    stockList.value = rows
  }
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
    refreshAllStockDetails()
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
  Object.assign(newInboundItem.value, { sku: '', name: '', description: '', uom: '', price: 0, cost: 0 })
}
const outboundChartData = computed(() => {
  const byDate = {}
  for (const o of outboundList.value) {
    const date = (new Date(o.date)).toISOString().slice(0, 10)
    byDate[date] = (byDate[date] || 0) + o.quantity
  }
  const dates = Object.keys(byDate).sort()
  return {
    labels: dates.map(d => d.split('-').reverse().join('.')),
    datasets: [{
      label: "Отгрузки",
      data: dates.map(d => byDate[d]),
      borderColor: "#ef4444",
      backgroundColor: "rgba(239,68,68,0.15)",
      tension: 0.3
    }]
  }
})

const filteredOutboundList = computed(() =>
  (outboundList.value || []).filter(o =>
    (!outboundDateFilter.value || (o.date || '').startsWith(outboundDateFilter.value)) &&
    (
      (o.name || '').toLowerCase().includes(outboundSearch.value.toLowerCase()) ||
      (o.sku || '').toLowerCase().includes(outboundSearch.value.toLowerCase()) ||
      (o.warehouse || '').toLowerCase().includes(outboundSearch.value.toLowerCase()) ||
      (o.destination || '').toLowerCase().includes(outboundSearch.value.toLowerCase())
    )
  )
)


const totalOutboundQuantity = computed(() => (outboundList.value || []).reduce((acc, o) => acc + o.quantity, 0))// Эту функцию вызывай при сохранении поставки
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
  // Если выбран новый товар
  if (newInbound.value.item_id === -1) {
    if (!newInboundItem.value.sku || !newInboundItem.value.name) {
      alert('Введите SKU и наименование товара');
      return;
    }
    const inboundPayload = {
      item_id: 0,
      supplier_id: newInbound.value.supplier_id,
      warehouse_id: newInbound.value.warehouse_id,
      quantity: newInbound.value.quantity,
      received_at: receivedAt,
    }
    const itemPayload = {
      sku: newInboundItem.value.sku,
      name: newInboundItem.value.name,
      description: newInboundItem.value.description,
      uom: newInboundItem.value.uom || 'шт',
      reorder_level: 0,
      reorder_qty: 0,
      price: newInboundItem.value.price,
      cost: newInboundItem.value.cost,
    }
    window.go.app.App.AddInboundTx(inboundPayload, itemPayload).then(() => {
      closeAddDeliveryModal()
      GetInboundDetails().then(data => { deliveriesList.value = data || [] })
      GetItems().then(data => { items.value = data || [] })
      refreshAllStockDetails()
    }).catch(err => {
      alert('Ошибка при добавлении поставки')
      console.error(err)
    })
  } else {
    // Если дата не выбрана, на бэке ставится now()
    const payload = {
      item_id: newInbound.value.item_id,
      supplier_id: newInbound.value.supplier_id,
      warehouse_id: newInbound.value.warehouse_id,
      quantity: newInbound.value.quantity,
      received_at: receivedAt,
      received_by: 1
    }
    window.go.app.App.AddInbound(payload).then(() => {
      closeAddDeliveryModal()
      // обновить deliveriesList после добавления
      GetInboundDetails().then(data => {
        deliveriesList.value = data || []
      })
      refreshAllStockDetails()
    }).catch(err => {
      alert("Ошибка при добавлении поставки")
      console.error(err)
    })
  }
}
async function reloadOutbound() {
  const data = await GetOutboundDetails();
  outboundList.value = Array.isArray(data) ? data : [];
}


watch(outboundDateFilter, reloadOutbound)

// Добавление
function openAddOutboundModal() {
  showAddOutboundModal.value = true
  newOutbound.value = { item_id: 0, warehouse_id: 0, destination: '', quantity: 1, shipped_at: '' }
}
function closeAddOutboundModal() { showAddOutboundModal.value = false }
async function confirmAddOutbound() {
  if (!newOutbound.value.item_id || !newOutbound.value.warehouse_id || !newOutbound.value.destination || !newOutbound.value.quantity) {
    alert('Заполните все поля')
    return
  }
  try {
    await AddOutbound(
  newOutbound.value.item_id,
  newOutbound.value.quantity,
  newOutbound.value.shipped_at,   // строка YYYY-MM-DD
  newOutbound.value.destination,
  newOutbound.value.warehouse_id
)
    closeAddOutboundModal()
    await  reloadOutbound()
    await reloadMovements()
    await refreshAllStockDetails()
  } catch (e) {
    alert('Ошибка при добавлении: ' + (e?.message || ''))
  }
}

// Редактирование
function openEditOutboundModal(o) {
  outboundToEdit.value = { ...o }
  showEditOutboundModal.value = true
}
function closeEditOutboundModal() { showEditOutboundModal.value = false }
async function confirmEditOutbound() {
  // в вашем outboundToEdit.value уже есть все поля:
  const o = outboundToEdit.value;

  // валидация...
  try {
    // передаём 6 отдельных параметров в том же порядке,
    // в каком вы объявили метод в app.go
    await EditOutbound(
      o.item_id,         // 1) itemID
      o.quantity,        // 2) quantity
      o.shipped_at,      // 3) shippedAtStr (YYYY-MM-DD)
      o.destination,     // 4) destination
      o.warehouse_id,    // 5) warehouseID
      o.outbound_id      // 6) outboundID
    );

    closeEditOutboundModal();
    await reloadOutbound();
    await reloadMovements();
    await refreshAllStockDetails();
  } catch (e) {
    alert('Ошибка при сохранении: ' + e?.message);
  }
}

// Удаление
async function deleteOutbound(o) {
  if (!confirm(`Удалить отгрузку товара "${o.name}"?`)) return
  try {
    await RemoveOutbound(o.outbound_id)
    reloadOutbound()
    refreshAllStockDetails()
  } catch (e) {
    alert('Ошибка при удалении: ' + (e?.message || ''))
  }
}
const weeklyStockChartData = computed(() => ({
  labels: weeklyStockData.value.map(d => formatDate(d.date)),
  datasets: [
    {
      label: 'Остатки',
      data: weeklyStockData.value.map(d => Number(d.total)),
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
      refreshAllStockDetails()
      const reload = selectedWarehouseId.value === 0
        ? GetStockDetails
        : () => FindStockByWarehouse(selectedWarehouseId.value)

      reload().then(data => {
        stockList.value = normalizeStockRows(data)
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
        refreshAllStockDetails()
      })
      .catch(err => {
        alert("Ошибка при удалении поставки");
        console.error(err);
      });
  }
}

function exportDeliveriesToExcel() {
  ExportDeliveriesToExcel().then(base64data => {
    const binary = atob(base64data)
    const len = binary.length
    const bytes = new Uint8Array(len)
    for (let i = 0; i < len; i++) {
      bytes[i] = binary.charCodeAt(i)
    }
    const blob = new Blob([bytes], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = 'deliveries.xlsx'
    link.click()
    setTimeout(() => URL.revokeObjectURL(link.href), 1000)
  }).catch(err => {
    alert('Ошибка экспорта: ' + err)
  })
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
      reloadWeeklyTrend()
      refreshAllStockDetails()
      if (selectedWarehouseId.value === 0) {
        GetStockDetails().then(data => {
  stockList.value = normalizeStockRows(data)
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
  ExportItemsToExcel().then(base64data => {
    const binary = atob(base64data)
    const len = binary.length
    const bytes = new Uint8Array(len)
    for (let i = 0; i < len; i++) {
      bytes[i] = binary.charCodeAt(i)
    }
    const blob = new Blob([bytes], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = 'items.xlsx'
    link.click()
    setTimeout(() => URL.revokeObjectURL(link.href), 1000)
  }).catch(err => {
    alert('Ошибка экспорта: ' + err)
  })
}
function closeAddModal() {
  showAddModal.value = false
  newStock.value = { item_id: 0, warehouse_id: 0, quantity: 0 }
}
async function reloadWeeklyTrend() {
  weeklyStockData.value = await GetWeeklyStockTrend() || []
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
    reloadWeeklyTrend()
    refreshAllStockDetails()
    if (selectedWarehouseId.value === 0) {
      GetStockDetails().then(data => {
        stockList.value = normalizeStockRows(data)
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

async function addWarehouse() {
  if (!newWarehouse.value.name.trim()) return
  await AddWarehouse(newWarehouse.value)
  await loadWarehouses()
  closeAddModal()
}

async function updateWarehouse() {
  if (!editWarehouseData.value.name.trim()) return
  await EditWarehouse(editWarehouseData.value)
  await loadWarehouses()
  closeEditModal()
}

function editWarehouse(w) {
  editWarehouseData.value = { ...w }
  showEditModal.value = true
}

const filteredWarehouses = computed(() =>
  warehouses.value.filter(w =>
    w.name.toLowerCase().includes(warehouseSearch.value.toLowerCase()) ||
    (w.location || '').toLowerCase().includes(warehouseSearch.value.toLowerCase())
  )
)

const dashboardStockRows = computed(() => {
  const warehouseId = Number(selectedRiskWarehouseId.value)
  return allStockList.value.filter(row => warehouseId === 0 || row.warehouse_id === warehouseId)
})

const lowStockItems = computed(() => {
  const stockMap = new Map()
  for (const row of dashboardStockRows.value) {
    const current = stockMap.get(row.item_id)
    if (!current) {
      stockMap.set(row.item_id, {
        warehouse_id: row.warehouse_id,
        warehouse: row.warehouse,
        quantity: Number(row.quantity) || 0
      })
      continue
    }
    current.quantity += Number(row.quantity) || 0
  }

  return items.value
    .map(item => {
      const stockInfo = stockMap.get(item.item_id)
      const currentStock = stockInfo?.quantity ?? 0
      const reorderLevel = Number(item.reorder_level) || 0
      const reorderQty = Number(item.reorder_qty) || 0
      const shortage = Math.max(reorderLevel - currentStock, 0)
      const recommendedOrder = shortage > 0
        ? Math.max(reorderQty, shortage)
        : 0

      if (reorderLevel <= 0 || currentStock > reorderLevel) {
        return null
      }

      const severity = currentStock === 0 ? 'critical' : currentStock <= Math.ceil(reorderLevel * 0.4) ? 'warning' : 'attention'

      return {
        item_id: item.item_id,
        warehouse_id: stockInfo?.warehouse_id || 0,
        warehouse: stockInfo?.warehouse || (selectedRiskWarehouseId.value === 0 ? 'Не распределен' : 'Нет остатков'),
        name: item.name,
        sku: item.sku,
        current_stock: currentStock,
        reorder_level: reorderLevel,
        recommended_order: recommendedOrder,
        estimated_cost: recommendedOrder * (Number(item.cost) || 0),
        severity,
        status_label: severity === 'critical' ? 'Критично' : severity === 'warning' ? 'Низкий остаток' : 'Нужен контроль'
      }
    })
    .filter(Boolean)
    .sort((a, b) => {
      if (a.severity !== b.severity) {
        return ['critical', 'warning', 'attention'].indexOf(a.severity) - ['critical', 'warning', 'attention'].indexOf(b.severity)
      }
      return a.current_stock - b.current_stock
    })
})

const criticalLowStockCount = computed(() =>
  lowStockItems.value.filter(item => item.severity === 'critical').length
)

const recommendedRestockUnits = computed(() =>
  lowStockItems.value.reduce((acc, item) => acc + item.recommended_order, 0)
)

const recommendedRestockCost = computed(() =>
  lowStockItems.value.reduce((acc, item) => acc + item.estimated_cost, 0)
)

const warehouseRiskSummary = computed(() =>
  warehouses.value
    .map(warehouse => {
      const scoped = lowStockItems.value.filter(item => item.warehouse_id === warehouse.warehouse_id)
      return {
        warehouse_id: warehouse.warehouse_id,
        name: warehouse.name,
        low_count: scoped.length,
        critical_count: scoped.filter(item => item.severity === 'critical').length
      }
    })
    .filter(item => item.low_count > 0)
    .sort((a, b) => {
      if (b.critical_count !== a.critical_count) return b.critical_count - a.critical_count
      return b.low_count - a.low_count
    })
)

function formatMoney(value) {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    maximumFractionDigits: 0
  }).format(Number(value) || 0)
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
  await GetWeeklyStockTrend().then(data => weeklyStockData.value = data)

  GetDashboard().then(data => {
    totalStock.value = data.total_stock
    itemCount.value = data.item_count
    monthlyOrders.value = data.monthly_orders
    newItems.value = data.new_items
  })
  GetSuppliers().then(data => suppliers.value = data || [])
  GetStockDetails().then(data => {
    const rows = normalizeStockRows(data)
    stockList.value = rows
    allStockList.value = rows
  })

  await reloadMovements()
  await reloadOutbound()
  GetInboundDetails().then(data => {
    deliveriesList.value = data || [];
  }).catch(err => {
    console.error("Ошибка загрузки поставок:", err)
  })
  GetTopItems().then(data => topItems.value = data)
  GetWarehouses().then(data => warehouses.value.push(...data))
  GetTurnoverByWarehouse().then(data => turnoverData.value = data)
})
watch([moveType, moveWarehouseId], reloadMovements)
watch(currentTab, async (tab) => {
  if (tab === 'Дашборд') {
    const data = await GetDashboard()
    totalStock.value = data.total_stock
    itemCount.value = data.item_count
    monthlyOrders.value = data.monthly_orders
    newItems.value = data.new_items
    if (tab === 'Движения') {
    await reloadMovements()
    await reloadOutbound()
    }
    topItems.value = await GetTopItems() || []
    turnoverData.value = await GetTurnoverByWarehouse() || []
    await refreshAllStockDetails()
  }
  if (tab === 'Пользователи' && user.value?.role === 'admin') {
    try {
      users.value = await GetUsers() || []
    } catch (e) {
      users.value = []
    }
  }
})


watch(selectedWarehouseId, (id) => {
  const warehouseId = Number(id)
  if (warehouseId === 0) {
    GetStockDetails().then(data => {
      const rows = normalizeStockRows(data)
      stockList.value = rows
      allStockList.value = rows
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
.dashboard-risk-grid {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(280px, 1fr);
  gap: 1.2rem;
  margin-top: 1.2rem;
}

.dashboard-risk-card,
.dashboard-risk-sidecard {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(247, 249, 252, 0.98));
}

.dashboard-risk-header {
  align-items: flex-start;
  gap: 1rem;
}

.dashboard-subtitle,
.dashboard-row-meta {
  color: #6b7280;
  font-size: 0.85rem;
}

.dashboard-risk-filter {
  min-width: 220px;
  max-width: 320px;
}

.risk-summary-cards {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.8rem;
  margin-bottom: 1rem;
}

.risk-summary-card {
  border-radius: 16px;
  padding: 0.9rem 1rem;
  border: 1px solid rgba(15, 23, 42, 0.08);
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.risk-summary-card strong {
  font-size: 1.15rem;
  color: #111827;
}

.risk-summary-card.warning { background: #fff7ed; }
.risk-summary-card.danger { background: #fef2f2; }
.risk-summary-card.neutral { background: #f8fafc; }
.risk-summary-card.accent { background: #eff6ff; }

.risk-summary-label {
  color: #4b5563;
  font-size: 0.84rem;
}

.stock-health-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.32rem 0.7rem;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 700;
  letter-spacing: 0.01em;
}

.stock-health-badge.critical {
  background: #fee2e2;
  color: #b91c1c;
}

.stock-health-badge.warning {
  background: #ffedd5;
  color: #c2410c;
}

.stock-health-badge.attention {
  background: #fef9c3;
  color: #a16207;
}

.warehouse-risk-list {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.warehouse-risk-row {
  display: flex;
  justify-content: space-between;
  gap: 0.75rem;
  align-items: center;
  padding: 0.9rem 1rem;
  border-radius: 16px;
  background: #f8fafc;
  border: 1px solid rgba(15, 23, 42, 0.08);
}

.warehouse-risk-name {
  font-weight: 700;
  color: #111827;
}

.warehouse-risk-values {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  color: #6b7280;
  font-size: 0.82rem;
}

.warehouse-risk-values strong {
  color: #b91c1c;
  font-size: 1.1rem;
}

@media (max-width: 1100px) {
  .dashboard-risk-grid {
    grid-template-columns: 1fr;
  }

  .risk-summary-cards {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .dashboard-risk-filter {
    min-width: 100%;
    max-width: 100%;
  }

  .risk-summary-cards {
    grid-template-columns: 1fr;
  }
}
</style>

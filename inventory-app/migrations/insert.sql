INSERT INTO "User" (username, password, full_name, role) VALUES
('admin', 'adminpass', 'Администратор', 'admin'),
('manager1', 'managerpass', 'Пётр Менеджер', 'manager'),
('manager2', 'managerpass2', 'Мария Менеджер', 'manager'),
('worker1', 'workerpass1', 'Иван Рабочий', 'worker'),
('worker2', 'workerpass2', 'Анна Рабочая', 'worker');


INSERT INTO supplier (name, contact_info) VALUES
('ООО "Поставка+", Москва', 'info@postavka-plus.ru'),
('ЗАО "СнабЭксперт"', 'snab@example.com'),
('ИП "Доставкин"', 'deliveryman@sample.org');


INSERT INTO warehouse (name, location) VALUES
('Центральный склад', 'г. Москва'),
('Западный склад', 'г. Смоленск'),
('Резервный склад', 'г. Екатеринбург');

INSERT INTO inbound (item_id, supplier_id, quantity, received_by, warehouse_id) VALUES
(1, 1, 50, 1, 1),
(2, 2, 40, 2, 1),
(3, 1, 20, 3, 2),
(4, 2, 120, 2, 1),
(5, 3, 30, 1, 3);


INSERT INTO outbound (item_id, quantity, shipped_by, destination, warehouse_id) VALUES
(1, 10, 3, 'Сеть "Компик"', 1),
(2, 20, 2, 'Магазин "ТехноМир"', 1),
(3, 5, 4, 'Офис г. Казань', 2),
(5, 10, 2, 'Партнёр ZETA', 3);



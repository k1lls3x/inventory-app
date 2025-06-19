<h1 align="center"> :clipboard: Inventory App :clipboard: </h1>

<p align="center">
  <img src="https://img.shields.io/badge/Languages-2-blueviolet?style=for-the-badge">
  <img src="https://img.shields.io/github/repo-size/Zproger/bspwm-dotfiles?style=for-the-badge">
</p>

<p align="center">
  <a href="#-%D1%80%D1%83%D1%81%D1%81%D0%BA%D0%B0%D1%8F-%D0%B2%D0%B5%D1%80%D1%81%D0%B8%D1%8F"><img src="https://img.shields.io/badge/%D0%9F%D0%B5%D1%80%D0%B5%D0%B9%D1%82%D0%B8_%D0%BD%D0%B0_%D1%80%D1%83%D1%81%D1%81%D0%BA%D0%B8%D0%B9-red?style=for-the-badge"></a>
</p>

## :blue_book: About

<img src="screenshots/main.png" alt="main" align="right" width="500px">

Inventory App is a desktop inventory management system powered by [**Wails**](https://wails.io/) and Vue. It lets you manage items, warehouses, suppliers and deliveries. Stock levels and movements are visualised with charts. The app supports role based authentication, Excel export and stores data in PostgreSQL.

## 🖼️ Gallery
Place additional screenshots in the `screenshots/` folder and they will appear here:
![gallery2](screenshots/2.png)
![gallery3](screenshots/3.png)
![gallery4](screenshots/4.png)
![gallery5](screenshots/5.png)
![gallery6](screenshots/6.png)

## :blue_book: Installation
1. Install **Go** (v1.21+) and **Node.js** (v16+).
2. Install the Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`.
3. Inside `frontend` run `npm install` to fetch web dependencies.
4. Build the application with `wails build`.
5. Run the generated binary from `build/bin`.

## :moyai: Contributing
Feel free to open issues or pull requests.

---

## 🇷🇺 Русская версия

**Inventory App** \- это настольная система учёта, созданная на основе [**Wails**](https://wails.io/) и Vue. Здесь можно управлять товарами, складами, поставщиками и отгрузками. Движение и количество товаров отображается на графиках. Есть ролевая авторизация, экспорт в **Excel** и хранение данных в PostgreSQL.

### Установка
1. Установите **Go** версии 1.21 и выше и **Node.js** версии 16 и выше.
2. Установите инструмент Wails командой `go install github.com/wailsapp/wails/v2/cmd/wails@latest`.
3. Перейдите в каталог `frontend` и выполните `npm install`.
4. Соберите программу командой `wails build`.
5. Запустите полученный бинарник из `build/bin`.

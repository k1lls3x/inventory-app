# :clipboard: Inventory App :clipboard:

[English version](README.md)

## :blue_book: О проекте

<img src="screenshots/main.png" alt="main" align="right" width="500px">

Inventory App — настольная система учёта, созданная на [**Wails**](https://wails.io/) и Vue. Она позволяет управлять товарами, складами, поставщиками и поставками, отслеживать остатки на складах с помощью интерактивных графиков и использовать авторизацию с ролями. Данные можно выгружать в **Excel**, всё хранится в PostgreSQL.

## 🖼️ Галерея
Дополнительные скриншоты можно разместить в папке `screenshots/` и они появятся здесь:
![gallery1](screenshots/1.png)
![gallery2](screenshots/2.png)
![gallery3](screenshots/3.png)
![gallery4](screenshots/4.png)
![gallery5](screenshots/5.png)
![gallery6](screenshots/6.png)

## :blue_book: Установка

### Требования
- Go 1.24 или новее
- Node.js 18 или новее
- PostgreSQL

### Шаги
1. Клонируйте репозиторий и перейдите в каталог проекта:
   ```bash
   git clone <repo-url>
   cd inventory-app-directory
   ```
2. Скопируйте `.env.example` в `.env` и задайте параметры подключения к базе данных.
3. Установите CLI Wails:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
4. Установите зависимости фронтенда:
   ```bash
   cd frontend
   npm install
   cd ..
   ```
5. Запустите миграции для PostgreSQL:
   ```bash
   psql -d <имя_бд> -f migrations/001_init_schema.sql
   ```
6. Запустите приложение в режиме разработки:
   ```bash
   wails dev
   ```
   или соберите production-версию:
   ```bash
   wails build
   ```
   Готовый бинарник появится в каталоге `build/bin`.

## :moyai: Как внести вклад
Присылайте issues и pull requests.


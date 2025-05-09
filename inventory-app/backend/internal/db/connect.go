package db

import (
	"os"
	"log"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)
var DB *sqlx.DB
func Init() {

	_ = godotenv.Load() // Загружаем переменные из .env
	cfg := LoadConfigFromEnv()
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		log.Println("❌ Ошибка подключения к PostgreSQL:", err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		log.Println("❌ PostgreSQL не отвечает:", err)
		os.Exit(1)
	}

	DB = db
	log.Println("✅ Подключение к PostgreSQL успешно")
}

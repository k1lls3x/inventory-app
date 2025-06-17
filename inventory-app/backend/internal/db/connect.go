package db

import (
	"os"
	"log"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Init() *sqlx.DB {
	_ = godotenv.Load()
	cfg := LoadConfig()
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
			log.Println("❌ Ошибка подключения к PostgreSQL:", err)
			os.Exit(1)
	}

	if err := db.Ping(); err != nil {
			log.Println("❌ PostgreSQL не отвечает:", err)
			os.Exit(1)
	}

	log.Println("✅ Подключение к PostgreSQL успешно")
	return db
}

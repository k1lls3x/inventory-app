package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"inventory-app/migrations"
	"log"
	"os"
)

func Init() *sqlx.DB {
	loadDotEnv()
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

	if err := migrations.Apply(db); err != nil {
		log.Println("❌ Ошибка применения схемы:", err)
		os.Exit(1)
	}

	if err := ensureBootstrapAdmin(db); err != nil {
		log.Println("❌ Ошибка создания стартового администратора:", err)
		os.Exit(1)
	}

	log.Println("✅ Подключение к PostgreSQL успешно")
	return db
}

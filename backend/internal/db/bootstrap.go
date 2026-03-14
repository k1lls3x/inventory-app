package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultBootstrapAdminUsername = "admin"
	defaultBootstrapAdminPassword = "admin123"
	defaultBootstrapAdminFullName = "Administrator"
)

func ensureBootstrapAdmin(db *sqlx.DB) error {
	var userCount int
	if err := db.Get(&userCount, `SELECT COUNT(*) FROM "User"`); err != nil {
		return err
	}

	if userCount > 0 {
		return nil
	}

	username := firstNonEmpty(os.Getenv("BOOTSTRAP_ADMIN_USERNAME"), defaultBootstrapAdminUsername)
	password := firstNonEmpty(os.Getenv("BOOTSTRAP_ADMIN_PASSWORD"), defaultBootstrapAdminPassword)
	fullName := firstNonEmpty(os.Getenv("BOOTSTRAP_ADMIN_FULL_NAME"), defaultBootstrapAdminFullName)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if _, err := db.Exec(
		`INSERT INTO "User"(username, password_hash, full_name, role) VALUES ($1, $2, $3, 'admin')`,
		username,
		string(passwordHash),
		fullName,
	); err != nil {
		return err
	}

	log.Printf("Создан стартовый администратор %q", username)
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

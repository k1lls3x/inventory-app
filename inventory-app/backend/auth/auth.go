package auth

import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	query := `
		SELECT user_id, username, password AS password_hash, role, full_name
		FROM "User"
		WHERE username = $1
	`
	err := db.DB.Get(&user, query, username)
	if err != nil {
		log.Println("❌ Произошла ошибка при получении пользователя", err)
		return nil, err
	}
	return &user, nil
}

func Register(user *model.User, rawPassword string) error {
	hash, err := HashPassword(rawPassword)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec(`
		INSERT INTO "User"(username, password, role, full_name) VALUES ($1, $2, $3, $4)`,
		user.Username, hash, user.Role, user.FullName,
	)
	return err
}

func Authorize(login, password string) (*model.User, bool) {
	user, err := GetUserByUsername(login)
	if err != nil {
		return nil, false
	}
	if CheckPasswordHash(password, user.PasswordHash) {
		return user, true
	}
	return nil, false
}

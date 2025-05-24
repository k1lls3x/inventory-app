package repository

import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
		"inventory-app/backend/auth"
)


func RemoveUser(userId int) error {

	query := `DELETE FROM "User" WHERE user_id = $1`
	_, err := db.DB.Exec(query, userId)

	if err != nil {
		log.Println("❌ Произошла ошибка при удалении пользователя: ", err)
	}
	return err
}

func GetUsers()([]model.User, error){

	var users []model.User
	query := `SELECT user_id, username, full_name, role FROM "User"
`
	err := db.DB.Select(&users, query)

	if err != nil {
		log.Println("❌ Произошла ошибка при получении пользователей: ", err)
		return nil, err
	}
	return users, err
}

func ChangeUserData(u model.UserUpdate) error {
	if u.Password != "" {
		// Меняем пароль
		hash, err := auth.HashPassword(u.Password)
		if err != nil {
			return err
		}
		_, err = db.DB.Exec(`
			UPDATE "User" SET username=$1, full_name=$2, role=$3, password_hash=$4 WHERE user_id=$5`,
			u.Username, u.FullName, u.Role, hash, u.UserID,
		)
		return err
	} else {
		// Не меняем пароль
		_, err := db.DB.Exec(`
			UPDATE "User" SET username=$1, full_name=$2, role=$3 WHERE user_id=$4`,
			u.Username, u.FullName, u.Role, u.UserID,
		)
		return err
	}
}


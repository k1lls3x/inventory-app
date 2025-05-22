package repository

import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
)

func GetSuppliers() ([]model.Supplier, error){
		var result []model.Supplier
		query := `
	SELECT
		supplier_id, name, inn, contact_person, phone, email
	FROM supplier
	ORDER BY name
`
		err := db.DB.Select(&result, query)

		if err != nil {
			log.Println("❌ Ошибка при получении данных о поставщиках: ", err)
		}
		return result, err
}

func AddSupplier(sp model.Supplier) error {
	query:=`
	INSERT INTO supplier (name, inn, contact_person, phone, email)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := db.DB.Exec(query, sp.Name, sp.Inn, sp.ContactPerson, sp.Phone, sp.Email)
		if err != nil {
		log.Println("❌ Ошибка при добавление данных поставщика: ", err)
	}
	return err
}

func EditSupplier(sp model.Supplier) error {

	query :=`
        UPDATE supplier
        SET
            name = $1,
            inn = $2,
            contact_person = $3,
            phone = $4,
            email = $5
        WHERE supplier_id = $6
    `
    _, err := db.DB.Exec(query, sp.Name, sp.Inn, sp.ContactPerson, sp.Phone, sp.Email, sp.SupplierID)
	if err != nil {
		log.Println("❌ Ошибка при изменение данных о поставщиках: ", err)
	}
	return err

}

func RemoveSupplier(supplierId int) error {
		query:=`
		DELETE FROM supplier WHERE supplier_id = $1;
		`
	_, err := db.DB.Exec(query, supplierId)
	if err != nil {
		log.Println("❌ Ошибка при удалении поставщика: ", err)
	}
	return err
}

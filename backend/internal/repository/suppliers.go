package repository

import (
	"context"
	"inventory-app/backend/internal/model"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type SupplierRepository interface {
	GetSuppliers(ctx context.Context) ([]model.Supplier, error)
	AddSupplier(ctx context.Context, sp model.Supplier) error
	EditSupplier(ctx context.Context, sp model.Supplier) error
	RemoveSupplier(ctx context.Context, supplierId int) error
}

type PgSupplierRepository struct {
	db  *sqlx.DB
	log zerolog.Logger
}

func NewSupplierRepository(db *sqlx.DB, log zerolog.Logger) *PgSupplierRepository {
	return &PgSupplierRepository{db: db, log: log}
}

func (r *PgSupplierRepository) GetSuppliers(ctx context.Context) ([]model.Supplier, error) {
	var result []model.Supplier
	query := `
	SELECT
		supplier_id, name, inn, contact_person, phone, email
	FROM supplier
	ORDER BY name`
	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении данных о поставщиках")
		return nil, err
	}
	r.log.Info().Int("count", len(result)).Msg("Список поставщиков успешно получен")
	return result, nil
}

func (r *PgSupplierRepository) AddSupplier(ctx context.Context, sp model.Supplier) error {
	query := `
	INSERT INTO supplier (name, inn, contact_person, phone, email)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, sp.Name, sp.Inn, sp.ContactPerson, sp.Phone, sp.Email)
	if err != nil {
		r.log.Error().
			Str("name", sp.Name).
			Err(err).Msg("Ошибка при добавлении поставщика")
		return err
	}
	r.log.Info().
		Str("name", sp.Name).
		Msg("Поставщик успешно добавлен")
	return nil
}

func (r *PgSupplierRepository) EditSupplier(ctx context.Context, sp model.Supplier) error {
	query := `
        UPDATE supplier
        SET
            name = $1,
            inn = $2,
            contact_person = $3,
            phone = $4,
            email = $5
        WHERE supplier_id = $6
    `
	_, err := r.db.ExecContext(ctx, query, sp.Name, sp.Inn, sp.ContactPerson, sp.Phone, sp.Email, sp.SupplierID)
	if err != nil {
		r.log.Error().
			Int("supplier_id", sp.SupplierID).
			Err(err).Msg("Ошибка при изменении данных поставщика")
		return err
	}
	r.log.Info().
		Int("supplier_id", sp.SupplierID).
		Msg("Данные поставщика успешно изменены")
	return nil
}

func (r *PgSupplierRepository) RemoveSupplier(ctx context.Context, supplierId int) error {
	query := `DELETE FROM supplier WHERE supplier_id = $1`
	res, err := r.db.ExecContext(ctx, query, supplierId)
	if err != nil {
		r.log.Error().
			Int("supplier_id", supplierId).
			Err(err).Msg("Ошибка при удалении поставщика")
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		r.log.Warn().Int("supplier_id", supplierId).Msg("Поставщик не найден для удаления")
	} else {
		r.log.Info().Int("supplier_id", supplierId).Msg("Поставщик успешно удалён")
	}
	return nil
}

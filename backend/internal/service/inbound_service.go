package service

import (
	"context"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type InboundService struct {
	repo     repository.InboundRepository
	itemRepo repository.ItemRepository
	log      zerolog.Logger
	db       *sqlx.DB
}

func NewInboundService(repo repository.InboundRepository,itemRepo repository.ItemRepository,  db *sqlx.DB, log zerolog.Logger) *InboundService {
	return &InboundService{
		repo: repo,
		itemRepo: itemRepo,
		log:  log,
		db:   db,

	}
}

func (s *InboundService) AddInboundTx(ctx context.Context, inb model.Inbound, item model.Item) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка открытия транзакции")
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	itemExists, err := s.itemRepo.ItemExistsTx(ctx, tx, item.SKU)

	if err != nil {
		tx.Rollback()
		return err
	}
	if !itemExists{
		if err := s.itemRepo.AddItemTx(ctx, tx, item); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err:= s.repo.AddInboundTx(ctx,tx,inb); err != nil{
		tx.Rollback()
		return err
	}
	if err:= tx.Commit();err != nil{
		s.log.Error().Err(err).Msg("Ошибка коммита транзакции")
		return err
	}
	return nil
}

func (s *InboundService) ListInboundDetails(ctx context.Context) ([]model.InboundDetails, error) {
	inbounds, err := s.repo.GetInboundDetails(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при получении списка поставок")
		return nil, err
	}
	s.log.Info().Int("count", len(inbounds)).Msg("Список поставок успешно получен")
	return inbounds, nil
}

func (s *InboundService) AddInbound(ctx context.Context, inb model.Inbound) error {
	err := s.repo.AddInbound(ctx, inb)
	if err != nil {
			s.log.Error().
					Int("item_id", inb.ItemID).
					Int("supplier_id", inb.SupplierID).
					Int("quantity", inb.Quantity).
					Str("received_at", inb.ReceivedAt.Format("2006-01-02")).
					Int("warehouse_id", inb.WarehouseID).
					Err(err).
					Msg("Ошибка при добавлении поставки (AddInboundSimple)")
			return err
	}
	s.log.Info().
			Int("item_id", inb.ItemID).
			Int("supplier_id", inb.SupplierID).
			Int("quantity", inb.Quantity).
			Str("received_at", inb.ReceivedAt.Format("2006-01-02")).
			Int("warehouse_id", inb.WarehouseID).
			Msg("Поставка успешно добавлена (AddInboundSimple)")
	return nil
}

func (s *InboundService) ListInboundDetailsByDate(ctx context.Context, date string) ([]model.InboundDetails, error) {
	inbounds, err := s.repo.GetInboundDetailsByDate(ctx, date)
	if err != nil {
		s.log.Error().Err(err).Str("date", date).Msg("Ошибка при получении поставок по дате")
		return nil, err
	}
	s.log.Info().Str("date", date).Int("count", len(inbounds)).Msg("Список поставок по дате успешно получен")
	return inbounds, nil
}

func (s *InboundService) DeleteInbound(ctx context.Context, inboundId int) error {
	err := s.repo.DeleteInbound(ctx, inboundId)
	if err != nil {
		s.log.Error().
			Int("inbound_id", inboundId).
			Err(err).
			Msg("Ошибка при удалении поставки")
		return err
	}
	s.log.Info().
		Int("inbound_id", inboundId).
		Msg("Поставка успешно удалена")
	return nil
}

func (s *InboundService) EditInbound(ctx context.Context, inb model.Inbound) error {
	err := s.repo.EditInbound(ctx, inb)
	if err != nil {
		s.log.Error().
			Int("inbound_id", inb.InboundID).
			Err(err).
			Msg("Ошибка при редактировании поставки")
		return err
	}
	s.log.Info().
		Int("inbound_id", inb.InboundID).
		Msg("Поставка успешно отредактирована")
	return nil
}


package service

import (
	"context"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type ItemService struct {
	repo repository.ItemRepository
	log  zerolog.Logger
}

func NewItemService(repo repository.ItemRepository, logger zerolog.Logger) *ItemService {
	return &ItemService{
		repo: repo,
		log:  logger,
	}
}

func (s *ItemService) ListItems(ctx context.Context) ([]model.Item, error) {
	items, err := s.repo.GetItems(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Не удалось получить список товаров")
		return nil, err
	}
	return items, nil
}

func (s *ItemService) ListItemsBriefs(ctx context.Context) ([]model.ItemBrief, error) {
	briefs, err := s.repo.GetItemBriefList(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Не удалось получить краткий список товаров")
		return nil, err
	}
	return briefs, nil
}

func (s *ItemService) AddItem(ctx context.Context, item model.Item) error {
	if err := s.repo.AddItem(ctx, item); err != nil {
		s.log.Error().Err(err).Str("sku", item.SKU).Msg("Ошибка добавления товара")
		return err
	}
	s.log.Info().Str("sku", item.SKU).Msg("Товар добавлен")
	return nil
}

func (s *ItemService) AddItemTx(ctx context.Context, tx *sqlx.Tx, item model.Item) error {
	if err := s.repo.AddItemTx(ctx, tx, item); err != nil {
		s.log.Error().Err(err).Str("sku", item.SKU).Msg("Ошибка добавления товара (Tx)")
		return err
	}
	s.log.Info().Str("sku", item.SKU).Msg("Товар добавлен (Tx)")
	return nil
}

func (s *ItemService) UpdateItem(ctx context.Context, item model.Item) error {
	if err := s.repo.UpdateItem(ctx, item); err != nil {
		s.log.Error().Err(err).Str("sku", item.SKU).Msg("Ошибка обновления товара")
		return err
	}
	s.log.Info().Str("sku", item.SKU).Msg("Товар обновлен")
	return nil
}

func (s *ItemService) RemoveItem(ctx context.Context, sku string) error {
	if err := s.repo.RemoveItem(ctx, sku); err != nil {
		s.log.Error().Err(err).Str("sku", sku).Msg("Ошибка удаления товара")
		return err
	}
	s.log.Info().Str("sku", sku).Msg("Товар удален")
	return nil
}

func (s *ItemService) FindItems(ctx context.Context, filter model.ItemFilter) ([]model.Item, error) {
	items, err := s.repo.FindItems(ctx, filter)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка поиска товаров")
		return nil, err
	}
	return items, nil
}


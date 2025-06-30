package Export

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
	"inventory-app/backend/internal/service"
)

func ExportItemsToExcel(itemService *service.ItemService, logger zerolog.Logger) (string, error) {
	items, err := itemService.ListItems(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка при получении товаров")
		return "", err
	}

	f := excelize.NewFile()
	sheet := "Товары"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"Наименование", "SKU", "Ед. изм.", "Описание", "Цена", "Себестоимость"}
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D9E1F2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	cellStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	f.SetColWidth(sheet, "A", "F", 22)

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for rowIdx, item := range items {
		price := ""
		if item.Price.Valid {
			price = fmt.Sprintf("%.2f", item.Price.Float64)
		}
		cost := ""
		if item.Cost.Valid {
			cost = fmt.Sprintf("%.2f", item.Cost.Float64)
		}
		desc := ""
		if item.Description != nil {
			desc = *item.Description
		}
		values := []interface{}{item.Name, item.SKU, item.UOM, desc, price, cost}
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, val)
			f.SetCellStyle(sheet, cell, cell, cellStyle)
		}
	}

	showStripes := true
	lastRow := len(items) + 1
	table := &excelize.Table{
		Range:          fmt.Sprintf("A1:F%d", lastRow),
		Name:           "ItemsTable",
		StyleName:      "TableStyleMedium9",
		ShowRowStripes: &showStripes,
	}

	if err := f.AddTable(sheet, table); err != nil {
		logger.Error().Err(err).Msg("Ошибка при добавлении таблицы")
		return "", err
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		logger.Error().Err(err).Msg("Ошибка при формировании файла")
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	logger.Info().Int("row_count", len(items)).Msg("Экспорт товаров в Excel завершён успешно")
	return encoded, nil
}

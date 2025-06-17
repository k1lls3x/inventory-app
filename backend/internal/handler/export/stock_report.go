package Export

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"context"
	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
	"inventory-app/backend/internal/service"
)


func ExportStockToExcel(stockService *service.StockService, logger zerolog.Logger) (string, error) {
	stockData, err := stockService.GetStockDetails(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка при получении остатков")
		return "", err
	}

	f := excelize.NewFile()
	sheet := "Остатки"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"Наименование", "Номер", "Склад", "Количество"}
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

	f.SetColWidth(sheet, "A", "D", 22)

	// Заголовки
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Данные
	for rowIdx, row := range stockData {
		values := []interface{}{row.Name, row.SKU, row.Warehouse, row.Quantity}
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, val)
			f.SetCellStyle(sheet, cell, cell, cellStyle)
		}
	}

	showStripes := true
	lastRow := len(stockData) + 1
	table := &excelize.Table{
		Range:          fmt.Sprintf("A1:D%d", lastRow),
		Name:           "StockTable",
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
	logger.Info().Int("row_count", len(stockData)).Msg("Экспорт остатков в Excel завершён успешно")
	return encoded, nil
}

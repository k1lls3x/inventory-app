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

func ExportSuppliersToExcel(supplierService *service.SupplierService, logger zerolog.Logger) (string, error) {
	suppliers, err := supplierService.GetSuppliers(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка при получении поставщиков")
		return "", err
	}

	f := excelize.NewFile()
	sheet := "Поставщики"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"Название", "ИНН", "Контактное лицо", "Телефон", "Email"}
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

	f.SetColWidth(sheet, "A", "E", 22)

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for rowIdx, s := range suppliers {
		inn := ""
		if s.Inn != nil {
			inn = *s.Inn
		}
		contact := ""
		if s.ContactPerson != nil {
			contact = *s.ContactPerson
		}
		phone := ""
		if s.Phone != nil {
			phone = *s.Phone
		}
		email := ""
		if s.Email != nil {
			email = *s.Email
		}
		values := []interface{}{s.Name, inn, contact, phone, email}
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, val)
			f.SetCellStyle(sheet, cell, cell, cellStyle)
		}
	}

	showStripes := true
	lastRow := len(suppliers) + 1
	table := &excelize.Table{
		Range:          fmt.Sprintf("A1:E%d", lastRow),
		Name:           "SuppliersTable",
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
	logger.Info().Int("row_count", len(suppliers)).Msg("Экспорт поставщиков в Excel завершён успешно")
	return encoded, nil
}

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

func ExportUsersToExcel(userService *service.UserService, logger zerolog.Logger) (string, error) {
	users, err := userService.GetUsers(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка при получении пользователей")
		return "", err
	}

	f := excelize.NewFile()
	sheet := "Пользователи"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "Логин", "ФИО", "Роль"}
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

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for rowIdx, u := range users {
		values := []interface{}{u.UserID, u.Username, u.FullName, u.Role}
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, val)
			f.SetCellStyle(sheet, cell, cell, cellStyle)
		}
	}

	showStripes := true
	lastRow := len(users) + 1
	table := &excelize.Table{
		Range:          fmt.Sprintf("A1:D%d", lastRow),
		Name:           "UsersTable",
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
	logger.Info().Int("row_count", len(users)).Msg("Экспорт пользователей в Excel завершён успешно")
	return encoded, nil
}

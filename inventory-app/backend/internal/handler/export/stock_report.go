package Export

import (
	"github.com/xuri/excelize/v2"
	"inventory-app/backend/internal/repository"
	"log"
	"encoding/base64"
	"bytes"
	"fmt"
)

func ExportStockToExcel() (string, error) {
	log.Println("Я родился ExportStockToExcel")
	stockData, err := repository.GetStockDetails()
	if err != nil {
			log.Println("❌ Ошибка при получении остатков:", err)
			return "", err
	}
	f := excelize.NewFile()
	sheet := "Остатки"
	index,_ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"Наименование", "Номер", "Склад", "Количество"}
	// Стили
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

	// Ширина
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
			Range:           fmt.Sprintf("A1:D%d", lastRow),
			Name:            "StockTable",
			StyleName:       "TableStyleMedium9",
			ShowRowStripes:  &showStripes,
	}

	err = f.AddTable(sheet, table)
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
			log.Println("❌ Ошибка при формировании файла:", err)
			return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return encoded, nil
}

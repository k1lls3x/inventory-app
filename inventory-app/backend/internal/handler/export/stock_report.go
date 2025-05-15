package export

import (
	"github.com/xuri/excelize/v2"
	"inventory-app/backend/internal/repository"
	"log"
)

func ExportStockToExcel(path string) error {
	log.Println("Я родился ExportStockToExcel")
		stockData, err := repository.GetStockDetails()

		if err != nil {
			log.Println("❌ Ошибка при получении остатков:", err)
			return err
		}
		f:= excelize.NewFile()
		sheet:= "Остатки"
		f.NewSheet(sheet)
		headers := []string{"Наименование","Номер","Склад","Количество"}
		for col, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			f.SetCellValue(sheet,cell,header)
		}
		for rowIdx, row := range stockData {
			values := []interface{}{row.Name, row.SKU, row.Warehouse, row.Quantity}
			for colIdx, val := range values {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
				f.SetCellValue(sheet, cell, val)
			}
		}
		if err := f.SaveAs(path); err != nil {
			log.Println("❌ Ошибка при сохранении файла:", err)
			return err
		}

		return nil
}


package excel

import (
	"fmt"

	"github.com/rmay1er/jgeo-excel/internal/models"
	"github.com/xuri/excelize/v2"
)

type ExcelWriter struct {
	file *excelize.File
}

func NewExcelWriter() *ExcelWriter {
	return &ExcelWriter{}
}

func (w *ExcelWriter) Write(data *[]models.CordsData, color ...string) error {
	if data == nil || len(*data) == 0 {
		return fmt.Errorf("нет данных для записи")
	}

	w.file = excelize.NewFile()
	sheetnum, err := w.file.NewSheet("geojson")
	if err != nil {
		return err
	}

	header := []any{"Тип", "Имя", "Описание", "Координаты"}

	if err := w.file.SetSheetRow(w.file.GetSheetName(sheetnum), "A1", header); err != nil {
		return err
	}

	for i, item := range *data {
		row := []any{item.Type, item.IconCaption, item.Description, item.Cords}
		cell := fmt.Sprintf("A%d", i+2)
		if err := w.file.SetSheetRow(w.file.GetSheetName(sheetnum), cell, row); err != nil {
			return err
		}
	}

	return nil
}

func (w *ExcelWriter) Save(path string) error {
	if err := w.file.SaveAs(path); err != nil {
		return err
	}
	return nil
}

func (w *ExcelWriter) Close() error {
	w.file.Close()
	w.file = nil
	return nil
}

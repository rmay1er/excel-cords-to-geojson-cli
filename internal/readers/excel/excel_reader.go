package readers

import (
	"fmt"

	"github.com/rmay1er/jgeo-excel/internal/models"
	"github.com/xuri/excelize/v2"
)

// ExcelReader читает координаты из Excel файла
type ExcelReader struct {
	file *excelize.File
	// Параметры для чтения
	sheet    string
	nameCol  string
	descCol  string
	cordsCol string
	startRow int
}

// NewExcelReader создает новый Excel reader
func NewExcelReader(path string, sheet, nameCol, descCol, cordsCol string, startRow int) (*ExcelReader, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть Excel файл: %w", err)
	}

	reader := &ExcelReader{
		file:     f,
		sheet:    sheet,
		nameCol:  nameCol,
		descCol:  descCol,
		cordsCol: cordsCol,
		startRow: startRow,
	}

	// Валидация параметров при создании
	if err := reader.validate(); err != nil {
		reader.Close()
		return nil, err
	}

	return reader, nil
}

// validate проверяет корректность параметров
func (r *ExcelReader) validate() error {
	// Проверяем, существует ли лист
	sheetIndex, err := r.file.GetSheetIndex(r.sheet)
	if err != nil || sheetIndex == -1 {
		return fmt.Errorf("лист '%s' не найден в файле", r.sheet)
	}

	// Валидируем колонки
	if r.descCol != "" {
		if _, err := excelize.ColumnNameToNumber(r.descCol); err != nil {
			return fmt.Errorf("неверное название колонки для описания: %v", err)
		}
	}

	if r.cordsCol != "" {
		if _, err := excelize.ColumnNameToNumber(r.cordsCol); err != nil {
			return fmt.Errorf("неверное название колонки для координат: %v", err)
		}
	}

	if r.nameCol != "" {
		if _, err := excelize.ColumnNameToNumber(r.nameCol); err != nil {
			return fmt.Errorf("неверное название колонки для названия: %v", err)
		}
	}

	return nil
}

// Read читает координаты из Excel файла
func (r *ExcelReader) Read() (*[]models.CordsData, error) {
	// Получаем все строки из листа
	rows, err := r.file.GetRows(r.sheet)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать строки из листа '%s': %w", r.sheet, err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("лист '%s' пуст", r.sheet)
	}

	// Конвертируем буквы колонок в индексы (A=1, B=2, C=3, ...)
	var nameColIdx int
	if r.nameCol != "" {
		nameColIdx, _ = excelize.ColumnNameToNumber(r.nameCol)
	}

	descColIdx, _ := excelize.ColumnNameToNumber(r.descCol)
	cordsColIdx, _ := excelize.ColumnNameToNumber(r.cordsCol)

	var result []models.CordsData

	// Начинаем с указанной строки (startRow обычно 2, т.к. 1я - заголовки)
	// startRow идет с 1, а индекс массива rows начинается с 0
	for i := r.startRow - 1; i < len(rows); i++ {
		row := rows[i]

		// Проверяем, что строка содержит координаты (это обязательное поле)
		if len(row) >= cordsColIdx && row[cordsColIdx-1] != "" {
			var cordsData models.CordsData

			// Добавляем координаты
			if err := cordsData.SetCords(row[cordsColIdx-1]); err != nil {
				// Пропускаем строку с ошибкой парсинга
				fmt.Printf("⚠️  Пропущена строка %d: ошибка при парсинге координат '%s'\n", i+1, row[cordsColIdx-1])
				continue
			}

			// Берем имя из соответствующей колонки (опционально, если колонка указана)
			if nameColIdx > 0 && len(row) >= nameColIdx && row[nameColIdx-1] != "" {
				cordsData.IconCaption = row[nameColIdx-1]
			}

			// Берем описание из соответствующей колонки
			if len(row) >= descColIdx && row[descColIdx-1] != "" {
				cordsData.Description = row[descColIdx-1]
			}

			// Добавляем объект в результат
			result = append(result, cordsData)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("не найдено координат в указанных колонках на листе '%s'", r.sheet)
	}

	return &result, nil
}

// Close закрывает Excel файл
func (r *ExcelReader) Close() error {
	return r.file.Close()
}

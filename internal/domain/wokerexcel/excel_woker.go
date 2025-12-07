package wokerexcel

import (
	"fmt"

	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/domain/cords"
	"github.com/xuri/excelize/v2"
)

type ExcelWoker struct {
	file *excelize.File
}

// NewExcelWoker создает новый воркер для работы с Excel файлом
func NewExcelWoker(path string) (*ExcelWoker, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	return &ExcelWoker{file: f}, nil
}

func (w *ExcelWoker) Close() error {
	return w.file.Close()
}

// GetCordsDataFromSheet читает координаты из указанного листа Excel
// Параметры:
// - sheet: название листа в Excel (например "Sheet1")
// - nameCol: буква колонки с названием (например "A"). Может быть пустой строкой, если название не требуется
// - descCol: буква колонки с описанием (например "B")
// - cordsCol: буква колонки с координатами (например "C")
// - startRow: номер строки, с которой начинать читать (обычно 2, т.к. 1я строка - заголовки)
func (ew *ExcelWoker) GetCordsDataFromSheet(sheet, nameCol, descCol, cordsCol string, startRow int) (*[]cords.CordsData, error) {
	// Проверяем, существует ли лист
	sheetIndex, err := ew.file.GetSheetIndex(sheet)
	if err != nil || sheetIndex == -1 {
		return nil, fmt.Errorf("лист '%s' не найден в файле", sheet)
	}

	// Смотрим все строчки на указанном листе
	rows, err := ew.file.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать строки из листа '%s': %w", sheet, err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("лист '%s' пуст", sheet)
	}

	// Конвертируем буквы колонок в индексы (A=1, B=2, C=3, ...)
	var nameColIdx int
	if nameCol != "" {
		var err error
		nameColIdx, err = excelize.ColumnNameToNumber(nameCol)
		if err != nil {
			return nil, fmt.Errorf("неверное название колонки для имени: %v", err)
		}
	}

	descColIdx, err := excelize.ColumnNameToNumber(descCol)
	if err != nil {
		return nil, fmt.Errorf("неверное название колонки для описания: %v", err)
	}
	cordsColIdx, err := excelize.ColumnNameToNumber(cordsCol)
	if err != nil {
		return nil, fmt.Errorf("неверное название колонки для координат: %v", err)
	}

	if descColIdx == 0 || cordsColIdx == 0 {
		return nil, fmt.Errorf("неверные названия колонок для описания и координат")
	}

	var result []cords.CordsData

	// Начинаем с указанной строки (startRow обычно 2, т.к. 1я - заголовки)
	// startRow идет с 1, а индекс массива rows начинается с 0
	for i := startRow - 1; i < len(rows); i++ {
		row := rows[i]

		// Проверяем, что строка содержит координаты (это обязательное поле)
		if len(row) >= cordsColIdx && row[cordsColIdx-1] != "" {
			var cordsData cords.CordsData

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
				cordsData.PointDesc = row[descColIdx-1]
			}

			// Добавляем объект в результат
			result = append(result, cordsData)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("не найдено координат в указанных колонках на листе '%s'", sheet)
	}

	return &result, nil
}

// GetCordsData старая функция для обратной совместимости
// Функция принимает колонку (например A или D) с именем будущей точки,
// колонку с описанием будущей точки и колонку с координатами точки.
// Использует первый лист и начинает со строки 2
func (ew *ExcelWoker) GetCordsData(nameCol, descCol, cordsCol string) (*[]cords.CordsData, error) {
	sheet := ew.file.GetSheetName(0)
	return ew.GetCordsDataFromSheet(sheet, nameCol, descCol, cordsCol, 2)
}

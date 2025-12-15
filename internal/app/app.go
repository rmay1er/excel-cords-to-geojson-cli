package app

import (
	"fmt"

	"github.com/rmay1er/jgeo-excel/internal/config"
	"github.com/rmay1er/jgeo-excel/internal/processors"

	xlsx "github.com/rmay1er/jgeo-excel/internal/readers/excel"
	"github.com/rmay1er/jgeo-excel/internal/writers"
	gjs "github.com/rmay1er/jgeo-excel/internal/writers/geojson"
)

// App основное приложение - фасад для работы с процессором
type JGeoApp struct {
	processor *processors.MarksProcessor
	writer    writers.Writer
	config    *config.Config
}

// NewApp создает новое приложение с процессором
func NewJGeoApp(processor *processors.MarksProcessor, writer writers.Writer) *JGeoApp {
	return &JGeoApp{
		processor: processor,
		writer:    writer,
	}
}

// NewAppWithConfig создает новое приложение с конфигурацией
func NewJGeoAppWithConfig(cfg *config.Config) (*JGeoApp, error) {
	// Создаем Reader для Excel
	excelReader, err := xlsx.NewExcelReader(
		cfg.Excel.File,
		cfg.Excel.Sheet,
		cfg.Excel.Columns.Name,
		cfg.Excel.Columns.Description,
		cfg.Excel.Columns.Coordinates,
		cfg.Excel.StartRow,
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать Excel reader: %w", err)
	}

	// Создаем Writer для GeoJSON
	geojsonWriter, err := gjs.NewGeojsonWriter(cfg.Geojson.Input)
	if err != nil {
		excelReader.Close()
		return nil, fmt.Errorf("не удалось создать GeoJSON writer: %w", err)
	}

	// Создаем процессор
	processor := processors.NewMarksProcessor(excelReader, geojsonWriter)

	return &JGeoApp{
		processor: processor,
		writer:    geojsonWriter,
		config:    cfg,
	}, nil
}

// Process выполняет основной процесс обработки координат
func (a *JGeoApp) ProcessToGeojson() error {
	if a.config == nil {
		return fmt.Errorf("конфигурация не установлена")
	}

	// Выполняем процесс обработки через процессор
	if err := a.processor.Process(a.config.Appearance.MarkerColor); err != nil {
		return err
	}

	// Сохраняем результат
	fmt.Printf("💾 Сохраняю результат в: %s\n", a.config.Geojson.Output)
	if err := a.writer.Save(a.config.Geojson.Output); err != nil {
		return fmt.Errorf("ошибка при сохранении GeoJSON файла: %w", err)
	}

	return nil
}

func (a *JGeoApp) ProcessToExcel(path string) error {

	// Выполняем процесс обработки через процессор
	if err := a.processor.Process(""); err != nil {
		return err
	}

	// Сохраняем результат
	fmt.Printf("💾 Сохраняю результат в: %s\n", path)
	if err := a.writer.Save(path); err != nil {
		return fmt.Errorf("ошибка при сохранении файла: %w", err)
	}

	return nil
}

// Close закрывает процессор и writer
func (a *JGeoApp) Close() error {
	if a.processor != nil {
		if err := a.processor.Close(); err != nil {
			return err
		}
	}
	return nil
}

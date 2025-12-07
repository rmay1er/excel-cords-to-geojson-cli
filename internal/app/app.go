package app

import (
	"fmt"

	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/config"
	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/processors"
	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/readers"
	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/writers"
)

// App основное приложение - фасад для работы с процессором
type App struct {
	processor *processors.CoordinatesProcessor
	writer    writers.Writer
	config    *config.Config
}

// NewApp создает новое приложение с процессором
func NewApp(processor *processors.CoordinatesProcessor, writer writers.Writer) *App {
	return &App{
		processor: processor,
		writer:    writer,
	}
}

// NewAppWithConfig создает новое приложение с конфигурацией
func NewAppWithConfig(cfg *config.Config) (*App, error) {
	// Создаем Reader для Excel
	excelReader, err := readers.NewExcelReader(
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
	geojsonWriter, err := writers.NewGeojsonWriter(cfg.Geojson.Input)
	if err != nil {
		excelReader.Close()
		return nil, fmt.Errorf("не удалось создать GeoJSON writer: %w", err)
	}

	// Создаем процессор
	processor := processors.NewCoordinatesProcessor(excelReader, geojsonWriter)

	return &App{
		processor: processor,
		writer:    geojsonWriter,
		config:    cfg,
	}, nil
}

// Process выполняет основной процесс обработки координат
func (a *App) Process() error {
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

// Close закрывает процессор и writer
func (a *App) Close() error {
	if a.processor != nil {
		if err := a.processor.Close(); err != nil {
			return err
		}
	}
	return nil
}

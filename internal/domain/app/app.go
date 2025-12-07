package app

import (
	"fmt"

	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/config"
	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/domain/wokerexcel"
	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/domain/wokergeo"
)

type App struct {
	excelWoker   *wokerexcel.ExcelWoker
	geojsonWoker *wokergeo.GeojsonWoker
	config       *config.Config
}

// NewApp создает новое приложение с воркерами
func NewApp(excelWoker *wokerexcel.ExcelWoker, geojsonWoker *wokergeo.GeojsonWoker) *App {
	return &App{
		excelWoker:   excelWoker,
		geojsonWoker: geojsonWoker,
	}
}

// NewAppWithConfig создает новое приложение с конфигурацией
func NewAppWithConfig(cfg *config.Config) (*App, error) {
	excelWoker, err := wokerexcel.NewExcelWoker(cfg.Excel.File)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть Excel файл: %w", err)
	}

	geojsonWoker, err := wokergeo.NewGeojsonWoker(cfg.Geojson.Input)
	if err != nil {
		excelWoker.Close()
		return nil, fmt.Errorf("не удалось открыть GeoJSON файл: %w", err)
	}

	return &App{
		excelWoker:   excelWoker,
		geojsonWoker: geojsonWoker,
		config:       cfg,
	}, nil
}

// Process обрабатывает координаты из Excel и добавляет их в GeoJSON
func (a *App) Process() error {
	if a.config == nil {
		return fmt.Errorf("конфигурация не установлена")
	}

	// Читаем координаты из Excel с указанным листом и начальной строкой
	cordsData, err := a.excelWoker.GetCordsDataFromSheet(
		a.config.Excel.Sheet,
		a.config.Excel.Columns.Name,
		a.config.Excel.Columns.Description,
		a.config.Excel.Columns.Coordinates,
		a.config.Excel.StartRow,
	)
	if err != nil {
		return fmt.Errorf("ошибка при чтении координат из Excel: %w", err)
	}

	// Добавляем каждую точку в GeoJSON
	for _, cord := range *cordsData {
		color := a.config.Appearance.MarkerColor
		if color == "" {
			color = "#FF0000" // Цвет по умолчанию
		}

		if err := a.geojsonWoker.AddPoints(cord, color); err != nil {
			return fmt.Errorf("ошибка при добавлении точки в GeoJSON: %w", err)
		}
	}

	// Сохраняем результат
	if err := a.geojsonWoker.SaveToGeojson(a.config.Geojson.Output); err != nil {
		return fmt.Errorf("ошибка при сохранении GeoJSON файла: %w", err)
	}

	return nil
}

// Close закрывает все воркеры
func (a *App) Close() error {
	if a.excelWoker != nil {
		if err := a.excelWoker.Close(); err != nil {
			return err
		}
	}
	if a.geojsonWoker != nil {
		if err := a.geojsonWoker.Close(); err != nil {
			return err
		}
	}
	return nil
}

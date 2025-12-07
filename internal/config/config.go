package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ColumnMapping описывает маппинг столбцов Excel
type ColumnMapping struct {
	Name        string
	Description string
	Coordinates string
}

// ExcelConfig конфигурация для работы с Excel файлом
type ExcelConfig struct {
	File     string
	Sheet    string
	Columns  ColumnMapping
	StartRow int
}

// GeojsonConfig конфигурация для работы с GeoJSON файлом
type GeojsonConfig struct {
	Input  string
	Output string
}

// AppearanceConfig конфигурация внешнего вида маркеров
type AppearanceConfig struct {
	MarkerColor string
}

// Config основная структура конфигурации
type Config struct {
	Excel      ExcelConfig
	Geojson    GeojsonConfig
	Appearance AppearanceConfig
}

// LoadConfig загружает конфигурацию из файла используя Viper
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	// Устанавливаем путь и имя файла конфигурации
	v.SetConfigFile(path)

	// Читаем файл конфигурации
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл конфигурации: %w", err)
	}

	// Парсим конфигурацию в структуру
	config := &Config{}

	// Excel конфигурация
	config.Excel.File = v.GetString("excel.file")
	config.Excel.Sheet = v.GetString("excel.sheet")
	config.Excel.Columns.Name = v.GetString("excel.columns.name")
	config.Excel.Columns.Description = v.GetString("excel.columns.description")
	config.Excel.Columns.Coordinates = v.GetString("excel.columns.coordinates")
	config.Excel.StartRow = v.GetInt("excel.start_row")

	// GeoJSON конфигурация
	config.Geojson.Input = v.GetString("geojson.input")
	config.Geojson.Output = v.GetString("geojson.output")

	// Appearance конфигурация
	config.Appearance.MarkerColor = v.GetString("appearance.marker_color")

	// Валидация конфигурации
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate проверяет валидность конфигурации
func (c *Config) Validate() error {
	if c.Excel.File == "" {
		return fmt.Errorf("путь к Excel файлу не указан (excel.file)")
	}
	// if c.Excel.Columns.Name == "" {
	// 	return fmt.Errorf("столбец для названия не указан (excel.columns.name)")
	// }
	if c.Excel.Columns.Description == "" {
		return fmt.Errorf("столбец для описания не указан (excel.columns.description)")
	}
	if c.Excel.Columns.Coordinates == "" {
		return fmt.Errorf("столбец для координат не указан (excel.columns.coordinates)")
	}
	if c.Geojson.Input == "" {
		return fmt.Errorf("путь к входному GeoJSON файлу не указан (geojson.input)")
	}
	if c.Geojson.Output == "" {
		return fmt.Errorf("путь к выходному GeoJSON файлу не указан (geojson.output)")
	}

	// Если лист не указан, используем Sheet1 по умолчанию
	if c.Excel.Sheet == "" {
		c.Excel.Sheet = "Sheet1"
	}

	// Если startRow не указан, используем 2 (т.к. 1я строка - заголовки)
	if c.Excel.StartRow == 0 {
		c.Excel.StartRow = 2
	}

	// Если цвет маркера не указан, используем красный по умолчанию
	if c.Appearance.MarkerColor == "" {
		c.Appearance.MarkerColor = "#FF0000"
	}

	return nil
}

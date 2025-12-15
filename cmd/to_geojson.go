/*
Copyright © 2025 Ruslan Mayer
*/
package cmd

import (
	"fmt"

	"github.com/rmay1er/jgeo-excel/internal/app"
	"github.com/rmay1er/jgeo-excel/internal/config"
	"github.com/spf13/cobra"
)

// toGeoJsonCmd представляет команду to-geojson
var toGeoJsonCmd = &cobra.Command{
	Use:   "to-geojson",
	Short: "Преобразовать координаты из Excel в GeoJSON",
	Long: `Команда to-geojson читает координаты из Excel файла и добавляет их в GeoJSON файл.

Требуется конфигурационный файл YAML с указанием:
- Пути к Excel и GeoJSON файлам
- Названий листа и столбцов в Excel
- Пути для выходного GeoJSON файла

Пример конфига (config.yaml):
  excel:
    file: "data.xlsx"
    sheet: "Sheet1"
    columns:
      name: "A"
      description: "B"
      coordinates: "C"
    start_row: 2

  geojson:
    input: "base.geojson"
    output: "result.geojson"

  appearance:
    marker_color: "#FF0000"

Использование:
  excel-cords-to-geojson to-geojson --config config.yaml
  excel-cords-to-geojson to-geojson -c config.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Получаем путь к конфигурационному файлу из флага
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return fmt.Errorf("ошибка при получении флага --config: %w", err)
		}

		if configPath == "" {
			return fmt.Errorf("флаг --config обязателен. Используйте: to-geojson --config config.yaml")
		}

		fmt.Printf("📂 Загружаю конфигурацию из: %s\n", configPath)

		// Загружаем конфигурацию
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("❌ ошибка при загрузке конфигурации: %w", err)
		}

		fmt.Println("✅ Конфигурация загружена успешно")
		fmt.Printf("  📊 Excel файл: %s (лист: %s)\n", cfg.Excel.File, cfg.Excel.Sheet)
		fmt.Printf("  📍 Столбцы: название=%s, описание=%s, координаты=%s\n",
			cfg.Excel.Columns.Name, cfg.Excel.Columns.Description, cfg.Excel.Columns.Coordinates)
		fmt.Printf("  🗺️  GeoJSON: %s → %s\n", cfg.Geojson.Input, cfg.Geojson.Output)

		// Создаем приложение с конфигом
		// Создаем приложение с конфигом
		fmt.Println("\n🔧 Инициализирую приложение...")
		application, err := app.NewGeoAppWithConfig(cfg)
		if err != nil {
			return fmt.Errorf("❌ ошибка при инициализации приложения: %w", err)
		}
		defer application.Close()

		// Обрабатываем данные
		fmt.Println("\n🔄 Начинаю преобразование координат...")
		if err := application.Process(); err != nil {
			return fmt.Errorf("❌ ошибка при обработке: %w", err)
		}

		fmt.Printf("\n✅ Успешно! Результат сохранен в: %s\n", cfg.Geojson.Output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(toGeoJsonCmd)

	// Добавляем флаг для пути к конфигурационному файлу
	toGeoJsonCmd.Flags().StringP("config", "c", "", "Путь к конфигурационному YAML файлу (обязателен)")
	toGeoJsonCmd.MarkFlagRequired("config")
}

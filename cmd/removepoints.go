/*
Copyright © 2025 Ruslan Mayer
*/
package cmd

import (
	"fmt"

	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/writers"
	"github.com/spf13/cobra"
)

// removepointsCmd представляет команду removepoints
var removepointsCmd = &cobra.Command{
	Use:   "removepoints",
	Short: "Удалить все точки из GeoJSON файла",
	Long: `Удалить все точки (объекты) из GeoJSON файла, оставив пустую коллекцию объектов.

Пример:
	excel-cords-to-geojson removepoints --file путь/к/файлу.geojson`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			return fmt.Errorf("ошибка: требуется флаг --file")
		}

		fmt.Printf("🗑️  Удаляю все точки из: %s\n", filePath)

		// Создаем GeoJSON writer
		writer, err := writers.NewGeojsonWriter(filePath)
		if err != nil {
			return fmt.Errorf("❌ ошибка загрузки GeoJSON файла: %w", err)
		}
		defer writer.Close()

		// Очищаем все точки из файла
		if err := writer.RemoveAllPoints(); err != nil {
			return fmt.Errorf("❌ ошибка при удалении точек: %w", err)
		}

		// Сохраняем пустой GeoJSON файл
		if err := writer.Save(filePath); err != nil {
			return fmt.Errorf("❌ ошибка сохранения GeoJSON файла: %w", err)
		}

		fmt.Printf("✅ Все точки удалены из %s\n", filePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removepointsCmd)

	// Здесь вы определите флаги и настройки конфигурации.
	removepointsCmd.Flags().StringP("file", "f", "", "Путь к GeoJSON файлу")
	removepointsCmd.MarkFlagRequired("file")

	// Cobra поддерживает Persistent Flags, которые будут работать для этой команды
	// и всех подкоманд, например:
	// removepointsCmd.PersistentFlags().String("foo", "", "Справка для foo")

	// Cobra поддерживает локальные флаги, которые будут работать только при вызове этой команды
	// напрямую, например:
	// removepointsCmd.Flags().BoolP("toggle", "t", false, "Справка для toggle")
}

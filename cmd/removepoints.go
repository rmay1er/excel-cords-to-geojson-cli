/*
Copyright © 2025 ИМЯ <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/domain/wokergeo"
	"github.com/spf13/cobra"
)

// removepointsCmd представляет команду removepoints
var removepointsCmd = &cobra.Command{
	Use:   "removepoints",
	Short: "Удалить все точки из GeoJSON файла",
	Long: `Удалить все точки (объекты) из GeoJSON файла, оставив пустую коллекцию объектов.

Пример:
		excel-cords-to-geojson removepoints --file путь/к/файлу.geojson`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			fmt.Println("Ошибка: требуется флаг --file")
			return
		}

		// Создать нового обработчика GeoJSON
		worker, err := wokergeo.NewGeojsonWoker(filePath)
		if err != nil {
			fmt.Printf("Ошибка загрузки GeoJSON файла: %v\n", err)
			return
		}

		// Очистить все объекты из файла
		worker.RemoveAllPoints()

		// Сохранить пустой GeoJSON файл
		if err := worker.SaveToGeojson(filePath); err != nil {
			fmt.Printf("Ошибка сохранения GeoJSON файла: %v\n", err)
			return
		}

		fmt.Printf("Все точки удалены из %s\n", filePath)
	},
}

func init() {
	rootCmd.AddCommand(removepointsCmd)

	// Здесь вы определите флаги и настройки конфигурации.
	removepointsCmd.Flags().StringP("file", "f", "", "Путь к GeoJSON файлу")

	// Cobra поддерживает Persistent Flags, которые будут работать для этой команды
	// и всех подкоманд, например:
	// removepointsCmd.PersistentFlags().String("foo", "", "Справка для foo")

	// Cobra поддерживает локальные флаги, которые будут работать только при вызове этой команды
	// напрямую, например:
	// removepointsCmd.Flags().BoolP("toggle", "t", false, "Справка для toggle")
}

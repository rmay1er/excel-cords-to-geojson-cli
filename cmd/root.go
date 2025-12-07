/*
Copyright © 2025 Ruslan Mayer
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd представляет базовую команду приложения
var rootCmd = &cobra.Command{
	Use:   "excel-cords-to-geojson",
	Short: "Преобразование координат из Excel в GeoJSON",
	Long: `excel-cords-to-geojson - это CLI инструмент для преобразования координат из Excel файла в GeoJSON формат.

Приложение позволяет:
  • Читать координаты и данные из Excel файла
  • Добавлять точки в существующий GeoJSON файл (FeatureCollection)
  • Конфигурировать столбцы Excel и параметры через YAML файл
  • Поддерживает несколько листов в Excel

Использование:
  excel-cords-to-geojson convert --config config.yaml

Для примера конфигурационного файла смотрите config.example.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

// Execute добавляет все подкоманды к корневой команде и устанавливает флаги.
// Это вызывается из main.main() один раз.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Здесь определяются глобальные флаги и конфигурация.
	// Cobra поддерживает persistent флаги, которые, если определены здесь,
	// будут глобальными для всего приложения.
}

/*
Copyright © 2025 Ruslan Mayer
*/
package cmd

import (
	"strings"

	"github.com/rmay1er/jgeo-excel/internal/app"
	"github.com/rmay1er/jgeo-excel/internal/processors"
	readers "github.com/rmay1er/jgeo-excel/internal/readers/geojson"
	writers "github.com/rmay1er/jgeo-excel/internal/writers/excel"
	"github.com/spf13/cobra"
)

// toExcelCmd представляет команду to-excel
var toExcelCmd = &cobra.Command{
	Use:   "to-excel",
	Short: "Преобразовать информацию из GeoJSON в Excel",
	Long:  `Команда to-excel читает коллекцию из GeoJSON файла и создаёт из них xlsx.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		in, _ := cmd.Flags().GetString("input")
		out, _ := cmd.Flags().GetString("output")
		if out == "" {
			out = strings.TrimSuffix(in, ".geojson") + ".xlsx"
		}
		geojsonReader, err := readers.NewGeoJSONReader(in)
		excelWriter := writers.NewExcelWriter()
		processor := processors.NewMarksProcessor(geojsonReader, excelWriter)
		if err != nil {
			return err
		}
		app := app.NewJGeoApp(processor, excelWriter)
		if err := app.ProcessToExcel(out); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(toExcelCmd)

	// Добавляем флаг для пути к конфигурационному файлу
	toExcelCmd.Flags().StringP("input", "i", "", "Путь к GeoJson файлу обязателен")
	toExcelCmd.MarkFlagRequired("input")
}

package writers

import (
	"fmt"
	"os"

	geojson "github.com/paulmach/go.geojson"
	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/models"
)

// GeojsonWriter пишет координаты в GeoJSON формат
type GeojsonWriter struct {
	file *geojson.FeatureCollection
}

// NewGeojsonWriter создает новый GeoJSON writer и загружает файл
func NewGeojsonWriter(path string) (*GeojsonWriter, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать GeoJSON файл: %w", err)
	}

	featureCollection, err := geojson.UnmarshalFeatureCollection(f)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить GeoJSON файл: %w", err)
	}

	return &GeojsonWriter{file: featureCollection}, nil
}

// Write добавляет координаты в GeoJSON коллекцию
func (w *GeojsonWriter) Write(data *[]models.CordsData, color string) error {
	if data == nil || len(*data) == 0 {
		return fmt.Errorf("нет данных для записи")
	}

	for _, cord := range *data {
		// Инвертируем координаты из Excel формата [широта, долгота] в GeoJSON формат [долгота, широта]
		coords := cord.Cords
		if len(coords) == 2 {
			// Координаты из Excel приходят в формате [широта, долгота], меняем на [долгота, широта]
			coords = []float64{coords[1], coords[0]}
		}

		newPoint := geojson.NewFeature(geojson.NewPointGeometry(coords))

		// Добавляем свойства
		if cord.IconCaption != "" {
			newPoint.SetProperty("iconCaption", cord.IconCaption)
		}
		if cord.PointDesc != "" {
			newPoint.SetProperty("description", cord.PointDesc)
		}
		if color != "" {
			newPoint.SetProperty("marker-color", color)
		}

		w.file.AddFeature(newPoint)
	}

	return nil
}

// RemoveAllPoints удаляет все точки (Point features) из коллекции
func (w *GeojsonWriter) RemoveAllPoints() error {
	if w.file == nil {
		return nil
	}

	var newFeatures []*geojson.Feature
	for _, feature := range w.file.Features {
		// Проверяем, является ли геометрия точкой
		if feature.Geometry == nil || feature.Geometry.Type != geojson.GeometryPoint {
			newFeatures = append(newFeatures, feature)
		}
	}
	w.file.Features = newFeatures
	return nil
}

// Save сохраняет GeoJSON в файл
func (w *GeojsonWriter) Save(path string) error {
	file, err := w.file.MarshalJSON()
	if err != nil {
		return fmt.Errorf("не удалось сериализовать GeoJSON: %w", err)
	}

	if err := os.WriteFile(path, file, 0644); err != nil {
		return fmt.Errorf("не удалось сохранить GeoJSON файл: %w", err)
	}

	return nil
}

// Close закрывает GeoJSON файл
func (w *GeojsonWriter) Close() error {
	w.file = nil
	return nil
}

package wokergeo

import (
	"os"

	geojson "github.com/paulmach/go.geojson"
	"github.com/rmay1er/excel-cords-to-geojson-cli/internal/domain/cords"
)

type GeojsonWoker struct {
	file *geojson.FeatureCollection
}

// Основной вокер. Содержит файл.
func NewGeojsonWoker(path string) (*GeojsonWoker, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	featureCollection, err := geojson.UnmarshalFeatureCollection(f)

	if err != nil {
		return nil, err
	}

	return &GeojsonWoker{file: featureCollection}, nil
}

// Добавляем точки с properties.
func (gj *GeojsonWoker) AddPoints(data cords.CordsData, color string) error {

	// Инвертируем координаты из Excel формата [широта, долгота] в GeoJSON формат [долгота, широта]
	coords := data.Cords
	if len(coords) == 2 {
		// Координаты из Excel приходят в формате [широта, долгота], меняем на [долгота, широта]
		coords = []float64{coords[1], coords[0]}
	}

	newPoint := geojson.NewFeature(geojson.NewPointGeometry(coords))

	// Добавляем свойства
	if data.IconCaption != "" {
		newPoint.SetProperty("iconCaption", data.IconCaption)
	}
	if data.PointDesc != "" {
		newPoint.SetProperty("description", data.PointDesc)
	}
	if color != "" {
		newPoint.SetProperty("marker-color", color)
	}

	gj.file.AddFeature(newPoint)
	return nil
}

// Удаляем все точки (Point features) из коллекции.
func (gj *GeojsonWoker) RemoveAllPoints() error {
	if gj.file == nil {
		return nil
	}

	var newFeatures []*geojson.Feature
	for _, feature := range gj.file.Features {
		// Проверяем, является ли геометрия точкой
		if feature.Geometry == nil || feature.Geometry.Type != geojson.GeometryPoint {
			newFeatures = append(newFeatures, feature)
		}
	}
	gj.file.Features = newFeatures
	return nil
}

// Сохраняем geoJSON файл.
func (gj *GeojsonWoker) SaveToGeojson(path string) error {
	file, err := gj.file.MarshalJSON()
	defer gj.Close()
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, file, 0644); err != nil {
		return err
	}
	return nil
}

// Закрываем файл.
func (w *GeojsonWoker) Close() error {
	w.file = nil
	return nil
}

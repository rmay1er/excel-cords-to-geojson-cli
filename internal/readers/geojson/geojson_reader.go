package geojson

import (
	"os"

	geojson "github.com/paulmach/go.geojson"
	"github.com/rmay1er/jgeo-excel/internal/models"
)

type GeoJSONReader struct {
	path string
}

func NewGeoJSONReader(path string) (*GeoJSONReader, error) {
	reader := &GeoJSONReader{
		path: path,
	}

	return reader, nil
}

func (r *GeoJSONReader) Read() (*[]models.CordsData, error) {
	parsed := []models.CordsData{}
	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, err
	}
	geoCollection, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return nil, err
	}

	for _, feture := range geoCollection.Features {
		name, ok := feture.Properties["iconCaption"].(string)
		if !ok {
			name = ""
		}
		desc, ok := feture.Properties["description"].(string)
		if !ok {
			desc = ""
		}

		var cords any
		switch feture.Geometry.Type {
		case geojson.GeometryPoint:
			cords = feture.Geometry.Point
		case geojson.GeometryLineString:
			cords = feture.Geometry.LineString
		case geojson.GeometryPolygon:
			cords = feture.Geometry.Polygon
		default:
			cords = nil
		}

		newFeture := models.CordsData{
			Type:        string(feture.Geometry.Type),
			IconCaption: name,
			Description: desc,
			Cords:       cords,
		}
		parsed = append(parsed, newFeture)
	}

	return &parsed, nil
}

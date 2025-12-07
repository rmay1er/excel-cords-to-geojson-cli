package writers

import "github.com/rmay1er/excel-cords-to-geojson-cli/internal/models"

// Writer определяет интерфейс для записи координат в различные форматы
type Writer interface {
	// Write записывает координаты в целевой формат
	Write(data *[]models.CordsData, color string) error
	// Save сохраняет данные в файл
	Save(path string) error
	// Close закрывает соединение с целевым файлом
	Close() error
}

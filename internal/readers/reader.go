package readers

import "github.com/rmay1er/excel-cords-to-geojson-cli/internal/models"

// Reader определяет интерфейс для чтения координат из различных источников
type Reader interface {
	// Read читает координаты из источника
	Read() (*[]models.CordsData, error)
	// Close закрывает соединение с источником
	Close() error
}

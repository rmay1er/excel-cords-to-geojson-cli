package processors

import (
	"fmt"

	"github.com/rmay1er/jgeo-excel/internal/readers"
	"github.com/rmay1er/jgeo-excel/internal/writers"
)

// MarkCoordinatesProcessor обрабатывает координаты: читает из Reader, пишет в Writer
type MarkCoordinatesProcessor struct {
	reader readers.Reader
	writer writers.Writer
}

// NewMarkCoordinatesProcessor создает новый процессор координат
func NewMarkCoordinatesProcessor(reader readers.Reader, writer writers.Writer) *MarkCoordinatesProcessor {
	return &MarkCoordinatesProcessor{
		reader: reader,
		writer: writer,
	}
}

// Process выполняет основной процесс: читает данные из Reader, пишет в Writer
func (p *MarkCoordinatesProcessor) Process(color string) error {
	// 1. Читаем данные из Reader
	fmt.Println("📖 Читаю данные из источника...")
	data, err := p.reader.Read()
	if err != nil {
		return fmt.Errorf("ошибка при чтении данных: %w", err)
	}
	fmt.Printf("✅ Прочитано %d координат\n", len(*data))

	// 2. Пишем данные через Writer
	fmt.Println("✍️  Записываю данные в целевой формат...")
	if err := p.writer.Write(data, color); err != nil {
		return fmt.Errorf("ошибка при записи данных: %w", err)
	}
	fmt.Println("✅ Данные записаны успешно")

	return nil
}

// Close закрывает Reader и Writer
func (p *MarkCoordinatesProcessor) Close() error {
	var firstErr error

	if p.reader != nil {
		if err := p.reader.Close(); err != nil {
			firstErr = err
		}
	}

	if p.writer != nil {
		if err := p.writer.Close(); err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

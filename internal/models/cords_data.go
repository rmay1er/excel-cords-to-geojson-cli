package models

import (
	"strconv"
	"strings"
)

type CordsDataType string
type CordsListType [][][]float64

const (
	Point   CordsDataType = "Point"
	Polygon CordsDataType = "Polygon"
)

type CordsData struct {
	Type        string
	IconCaption string
	Description string
	Cords       any
	Color       string
}

func (c *CordsData) SetCords(cords string) error {
	// Разделяем строку по запятым, точкам с запятой и пробелам, удаляем пустые элементы
	parts := strings.FieldsFunc(cords, func(r rune) bool {
		return r == ',' || r == ';' || r == ' '
	})
	floatCords := make([]float64, 0, len(parts))
	for _, str := range parts {
		if str == "" {
			continue
		}
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		floatCords = append(floatCords, val)
	}
	c.Cords = floatCords
	return nil
}

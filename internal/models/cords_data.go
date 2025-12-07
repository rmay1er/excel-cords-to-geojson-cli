package models

import (
	"strconv"
	"strings"
)

type CordsData struct {
	IconCaption string
	PointDesc   string
	Cords       []float64
	MarkerColor string
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

package hardware_service

import (
	"src/models"
	"strings"
)

type Hardware struct {
	Id          int
	Name        string
	Type        string
	Description string
}

func (h *Hardware) IsTypeValid() bool {
	if strings.ToLower(h.Type) == "single-board computer" || strings.ToLower(h.Type) == "microcontroller unit" || strings.ToLower(h.Type) == "sensor" {
		return true
	}
	return false
}

func (h *Hardware) Add() {
	hardware := models.Hardware{
		Name:        h.Name,
		Type:        h.Type,
		Description: h.Description,
	}
	models.AddHardware(hardware)
}

func (h *Hardware) Delete() {
	models.DeleteHardware(h.Id)
}

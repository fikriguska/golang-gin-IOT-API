package hardware_service

import (
	"src/models"
	"strings"
)

type Hardware struct {
	models.Hardware
}

func (h *Hardware) IsTypeValid() bool {
	if strings.ToLower(h.Type) == "single-board computer" || strings.ToLower(h.Type) == "microcontroller unit" || strings.ToLower(h.Type) == "sensor" {
		return true
	}
	return false
}

func (h *Hardware) IsExist() bool {
	return models.IsHardwareExistById(h.Id)
}

func (h *Hardware) Add() {
	models.AddHardware(h.Hardware)
}

func (h *Hardware) Delete() {
	models.DeleteHardware(h.Id)
}

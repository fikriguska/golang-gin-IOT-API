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

func (h *Hardware) CheckHardwareType(type_ string) bool {
	switch type_ {
	case "sensor":
		return models.IsHardwareTypedSensorById(h.Id)
	case "node":
		return models.IsHardwareTypedNodeById(h.Id)
	}

	return false
}

func (h *Hardware) Add() {
	models.AddHardware(h.Hardware)
}

func (h *Hardware) GetAll() models.HardwareList {
	var list models.HardwareList
	list.Sensor = models.GetAllHardwareTypedSensor()
	list.Node = models.GetAllHardwareTypedNode()

	return list
}

func (h *Hardware) Get() interface{} {

	isSensor := models.IsHardwareTypedSensorById(h.Id)
	IsNode := models.IsHardwareTypedNodeById(h.Id)
	hw := models.GetHardwareById(h.Id)

	if isSensor {
		var hardware models.HardwareSensorGet
		hardware.Id = hw.Id
		hardware.Name = hw.Name
		hardware.Type = hw.Type
		hardware.Description = hw.Description

		sensor := models.GetSensorByHardwareId(h.Id)
		hardware.Sensor.Name = sensor.Name
		hardware.Sensor.Unit = sensor.Unit

		return hardware
	} else if IsNode {
		var hardware models.HardwareNodeGet
		hardware.Id = hw.Id
		hardware.Name = hw.Name
		hardware.Type = hw.Type
		hardware.Description = hw.Description

		node := models.GetNodeByHardwareId(h.Id)
		hardware.Node.Name = node.Name
		hardware.Node.Location = node.Location

		return hardware
	}
	return nil
}

func (h *Hardware) Update(hardware models.HardwareUpdate) {
	models.UpdateHardware(hardware, h.Id)
}

func (h *Hardware) Delete() error {
	return models.DeleteHardware(h.Id)
}

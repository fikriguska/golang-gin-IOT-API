package sensor_service

import (
	"src/models"
)

type Sensor struct {
	models.Sensor
}

func (s *Sensor) Add() {

	// check if there is a hardware
	if s.Id_hardware != -1 {
		models.AddSensor(s.Sensor)
	} else {
		models.AddSensorNoHardware(s.Sensor)
	}
}

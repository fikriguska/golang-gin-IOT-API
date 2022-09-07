package sensor_service

import (
	"src/models"
)

type Sensor struct {
	models.Sensor
}

func (s *Sensor) IsExistAndOwner(id_user int) (exist bool, owner bool) {
	exist = models.IsSensorExistById(s.Id)
	if !exist {
		return exist, false
	}
	owner = (models.GetUserIdBySensorId(s.Id) == id_user)
	return exist, owner
}

func (s *Sensor) Add() {

	// check if there is a hardware
	if s.Id_hardware != -1 {
		models.AddSensor(s.Sensor)
	} else {
		models.AddSensorNoHardware(s.Sensor)
	}
}

func (s *Sensor) Delete() {
	models.DeleteSensor(s.Id)
}

package sensor_service

import (
	"src/models"
	"src/service/cache_service"
)

type Sensor struct {
	models.Sensor
}

func (s *Sensor) IsExistAndOwner(id_user int) (exist bool, owner bool) {
	sensor_cached, found := cache_service.Get("sensor", s.Id)

	if !found {
		exist = models.IsSensorExistById(s.Id)
		if !exist {
			return exist, false
		} else {
			var sensor_cached models.CachedSensor
			sensor := models.GetSensorById(s.Id)
			sensor_cached.Sensor = sensor
			cache_service.Set("sensor", sensor.Id, sensor_cached)
		}
		owner = (models.GetUserIdBySensorId(s.Id) == id_user)
		return exist, owner
	} else {
		var userId int
		if sensor_cached.(models.CachedSensor).User.Id == 0 {
			userId = models.GetUserIdBySensorId(s.Id)
			sc := sensor_cached.(models.CachedSensor)
			sc.User.Id = userId
			cache_service.Set("sensor", s.Id, sc)
		} else {
			userId = sensor_cached.(models.CachedSensor).User.Id
		}
		return true, (userId == id_user)
	}
}

func (s *Sensor) Add() {

	// check if there is a hardware
	if s.Id_hardware != -1 {
		models.AddSensor(s.Sensor)
	} else {
		models.AddSensorNoHardware(s.Sensor)
	}
}

func (s *Sensor) GetAll(id_user int, is_admin bool) []models.SensorList {
	var sensors []models.SensorList
	if is_admin {
		sensors = models.GetAllSensor()
	} else {
		sensors = models.GetAllSensorByUserId(id_user)
	}
	return sensors
}

func (s *Sensor) Get() models.SensorGet {
	var resp models.SensorGet

	sensor := models.GetSensorById(s.Id)
	channels := models.GetChannelBySensorId(s.Id)

	resp.Id = sensor.Id
	resp.Name = sensor.Name
	resp.Unit = sensor.Unit

	for i, c := range channels {
		resp.Channel = append(resp.Channel, models.SensorChannelGet{})
		resp.Channel[i].Time = c.Time
		resp.Channel[i].Value = c.Value
	}

	return resp
}

func (s *Sensor) Update(sensor models.SensorUpdate) {
	models.UpdateSensor(sensor, s.Id)
}

func (s *Sensor) Delete() {
	models.DeleteSensor(s.Id)
}

package store

import (
	"fmt"

	"github.com/august-kuhfuss/hotdamn/domain"
)

type CreateSensorParams struct {
	ID   string
	Name string
}

type FindSensorsFilter struct {
	IDs      []string
	Names    []string
	IsActive bool
}

type FindSensorsParams struct {
	Filter *FindSensorsFilter
}

type SensorStore interface {
	CreateOrUpdateSensor(params *CreateSensorParams) error
	FindSensors(params *FindSensorsParams) ([]domain.Sensor, error)
}

type ErrSensorsNotFound struct {
	Filter *FindSensorsFilter
}

func (e ErrSensorsNotFound) Error() string {
	return fmt.Sprint("No sensors found with filter: ", e.Filter)
}

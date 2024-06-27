//go:build dev
// +build dev

package mock

import (
	"github.com/august-kuhfuss/hotdamn/domain"
	"github.com/august-kuhfuss/hotdamn/store"
)

func (s *mock) CreateOrUpdateSensor(params *store.CreateSensorParams) error {
	panic("not implemented") // TODO: Implement
}

func (s *mock) FindSensors(params *store.FindSensorsParams) ([]domain.Sensor, error) {
	panic("not implemented") // TODO: Implement
}

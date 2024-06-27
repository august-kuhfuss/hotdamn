//go:build dev
// +build dev

package mock

import (
	"github.com/august-kuhfuss/hotdamn/domain"
	"github.com/august-kuhfuss/hotdamn/store"
)

func (s *mock) CreateMeasurement(params *store.CreateMeasurementParams) error {
	panic("not implemented") // TODO: Implement
}

func (s *mock) FindMeasurements(params *store.FindMeasurementsParams) ([]domain.Measurement, error) {
	panic("not implemented") // TODO: Implement
}

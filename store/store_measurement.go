package store

import (
	"fmt"
	"time"

	"github.com/august-kuhfuss/hotdamn/domain"
)

type FindMeasurementsFilter struct {
	SensorIDs    []string
	TimeframeMin time.Time
	TimeframeMax time.Time
	ValueMin     float32
	ValueMax     float32
}

type CreateMeasurementParams struct {
	SensorID  string
	Timestamp time.Time
	ValueK    float32
}

type FindMeasurementsParams struct {
	Unit   domain.MeasurementUnit
	Filter *FindMeasurementsFilter
}

type MeasurementStore interface {
	CreateMeasurement(params *CreateMeasurementParams) error
	FindMeasurements(params *FindMeasurementsParams) ([]domain.Measurement, error)
}

type ErrMeasurementsNotFound struct {
	Filter *FindMeasurementsFilter
}

func (e ErrMeasurementsNotFound) Error() string {
	return fmt.Sprint("No measurements found with filter: ", e.Filter)
}

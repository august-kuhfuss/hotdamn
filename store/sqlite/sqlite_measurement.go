package sqlite

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/august-kuhfuss/hotdamn/domain"
	"github.com/august-kuhfuss/hotdamn/store"
)

func (s *sqlite) CreateMeasurement(params *store.CreateMeasurementParams) error {
	q := s.DB.
		Insert("measurements").
		Columns("sensor_id", "timestamp", "value_k").
		Values(params.SensorID, params.Timestamp, params.ValueK)

	if _, err := q.Exec(); err != nil {
		return err
	}
	return nil
}

func (s *sqlite) FindMeasurements(params *store.FindMeasurementsParams) ([]domain.Measurement, error) {
	data := make([]domain.Measurement, 0)

	valueCol := func() string {
		switch params.Unit {
		case domain.MeasurementUnitCelsius:
			return "(value_k - 273.15)"
		case domain.MeasurementUnitKelvin:
			return "value_k"
		case domain.MeasurementUnitFahrenheit:
			return "(value_k - 273.15) * 1.8 + 32"
		default:
			return "value_k"
		}
	}()

	q := s.DB.
		Select("sensor_id", "timestamp", valueCol).
		From("measurements").
		Where(sq.Or{
			sq.Eq{"sensor_id": params.Filter.SensorIDs},
			sq.GtOrEq{"timestamp": params.Filter.TimeframeMin},
			sq.LtOrEq{"timestamp": params.Filter.TimeframeMax},
			sq.GtOrEq{"value_k": params.Filter.ValueMin},
			sq.LtOrEq{"value_k": params.Filter.ValueMax},
		})

	rows, err := q.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m domain.Measurement
		if err := rows.Scan(&m.Sensor.ID, &m.Timestamp, &m.Value); err != nil {
			return nil, err
		}
		m.Unit = params.Unit
		data = append(data, m)
	}

	return data, nil
}

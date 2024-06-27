package sqlite

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/august-kuhfuss/hotdamn/domain"
	"github.com/august-kuhfuss/hotdamn/store"
)

func (s *sqlite) CreateOrUpdateSensor(params *store.CreateSensorParams) error {
	q := s.DB.
		Insert("sensors").Options("OR REPLACE").
		Columns("id", "name").
		Values(params.ID, params.Name)

	if _, err := q.Exec(); err != nil {
		return err
	}
	return nil
}

func (s *sqlite) FindSensors(params *store.FindSensorsParams) ([]domain.Sensor, error) {
	q := s.DB.
		Select("id", "name", "is_active").
		From("sensors").
		Where(sq.Or{
			sq.Eq{"id": params.Filter.IDs},
			sq.Eq{"name": params.Filter.Names},
			sq.Eq{"is_active": params.Filter.IsActive},
		})

	rows, err := q.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]domain.Sensor, 0)
	for rows.Next() {
		var s domain.Sensor
		if err := rows.Scan(&s.ID, &s.Name, &s.IsActive); err != nil {
			return nil, err
		}
		data = append(data, s)
	}

	return data, nil
}

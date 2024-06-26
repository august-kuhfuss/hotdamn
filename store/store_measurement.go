package store

import "github.com/august-kuhfuss/hotdamn/domain"

type MeasurementStore interface {
	CreateEntry(entry domain.MeasurementEntry) error
	GetEntries() ([]domain.MeasurementEntry, error)
}

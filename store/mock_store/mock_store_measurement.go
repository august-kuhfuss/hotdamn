//go:build dev
// +build dev

package mockstore

import (
	"fmt"
	"strconv"
	"time"

	"github.com/august-kuhfuss/hotdamn/domain"
)

func (s *mock) CreateEntry(entry domain.MeasurementEntry) error {
	panic("not implemented") // TODO: Implement
}

func (s *mock) GetEntries() ([]domain.MeasurementEntry, error) {
	entries := make([]domain.MeasurementEntry, 25)
	for i := range entries {
		entries[i] = domain.MeasurementEntry{
			Timestamp: s.F.Time().Time(time.Now()),
			Device: domain.MeasurementDevice{
				ID:   strconv.Itoa(i),
				Name: fmt.Sprintf("Device %d", i),
			},
			Value: s.F.Float32(2, -50, 50),
			Unit:  "C",
		}
	}

	return entries, nil
}

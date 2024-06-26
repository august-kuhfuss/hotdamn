package domain

import "time"

type MeasurementEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Device    MeasurementDevice `json:"device"`
	Value     float32           `json:"value"`
	Unit      string            `json:"unit"`
}

type MeasurementDevice struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

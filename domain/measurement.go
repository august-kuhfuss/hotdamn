package domain

import (
	"time"
)

type MeasurementUnit string

const (
	MeasurementUnitCelsius    MeasurementUnit = "C"
	MeasurementUnitFahrenheit MeasurementUnit = "F"
	MeasurementUnitKelvin     MeasurementUnit = "K"
)

func (u MeasurementUnit) String() string {
	return string(u)
}

type Measurement struct {
	Sensor    Sensor          `json:"sensor"`
	Timestamp time.Time       `json:"timestamp"`
	Value     float32         `json:"value"`
	Unit      MeasurementUnit `json:"unit"`
}

func ConvertCelsiusToKelvin(celsius float32) float32 {
	return celsius + 273.15
}

func ConvertCelsiusToFahrenheit(celsius float32) float32 {
	return celsius*9/5 + 32
}

func ConvertFahrenheitToKelvin(fahrenheit float32) float32 {
	return (fahrenheit-32)*5/9 + 273.15
}

func ConvertFahrenheitToCelsius(fahrenheit float32) float32 {
	return (fahrenheit - 32) * 5 / 9
}

func ConvertKelvinToCelsius(kelvin float32) float32 {
	return kelvin - 273.15
}

func ConvertKelvinToFahrenheit(kelvin float32) float32 {
	return (kelvin-273.15)*9/5 + 32
}

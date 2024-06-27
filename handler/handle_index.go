package handler

import (
	"net/http"

	"github.com/august-kuhfuss/hotdamn/domain"
	"github.com/august-kuhfuss/hotdamn/store"
	"github.com/labstack/echo/v4"
)

func init() {
	handler.GET("/", func(c echo.Context) error {
		cc := c.(*customContext)

		unit := domain.MeasurementUnitCelsius

		params := &store.FindMeasurementsParams{
			Unit:   unit,
			Filter: &store.FindMeasurementsFilter{},
		}

		data, err := cc.Store.FindMeasurements(params)
		if err != nil {
			return echo.ErrInternalServerError
		}

		for i := range data {
			s, err := cc.Store.FindSensors(&store.FindSensorsParams{Filter: &store.FindSensorsFilter{IDs: []string{data[i].Sensor.ID}}})
			if err != nil {
				return echo.ErrInternalServerError
			}

			data[i].Sensor = s[0]
		}

		return c.JSON(http.StatusOK, data)
	})
}

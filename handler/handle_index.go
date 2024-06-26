package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func init() {
	handler.GET("/", func(c echo.Context) error {
		cc := c.(*customContext)

		data, err := cc.Store.GetEntries()
		if err != nil {
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, data)
	})
}

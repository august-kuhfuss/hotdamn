package handler

import (
	"net/http"

	"github.com/august-kuhfuss/hotdamn/store"
	"github.com/labstack/echo/v4"
)

var handler = echo.New()

type customContext struct {
	echo.Context
	Store store.Store
}

func setStore(store store.Store) {
	handler.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &customContext{c, store}
			return next(cc)
		}
	})
}

func New(store store.Store) http.Handler {
	setStore(store)
	return handler
}

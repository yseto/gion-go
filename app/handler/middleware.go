package handler

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func CheckXHR() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			xrw := c.Request().Header.Get("X-Requested-With")
			if strings.ToLower(xrw) == "xmlhttprequest" {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}

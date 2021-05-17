package middlewares

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ParseId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return echo.ErrBadRequest
		}

		c.Set("id", id)
		return next(c)
	}
}

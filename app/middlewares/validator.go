package middlewares

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv CustomValidator) Validate(i interface{}) error {
	//if err := cv.Validator.Struct(i); err != nil {
	//	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	//}
	return nil
}

func ValidateRequest(p interface{}) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := c.Bind(p); err != nil {
				return echo.ErrBadRequest
			}
			if err := c.Validate(p); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			c.Set("body", p)
			return next(c)
		}
	}
}

package middlewares

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return echo.NewHTTPError(http.StatusInternalServerError, cv.Validator.Struct(i).Error())
}

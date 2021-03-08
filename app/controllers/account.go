package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

var GetAuthenticated = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var LinkAnonymously = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var UpdateAccount = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}

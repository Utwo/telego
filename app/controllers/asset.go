package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

var UploadAsset = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var DuplicateAsset = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var DeleteAsset = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}

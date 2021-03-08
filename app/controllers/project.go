package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"telego/app/validators"
)

var GetProjects = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var GetBySlugOfLoggedInUser = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var GetAssetStats = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var CheckIfSlugIsAvailable = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}

var CreateProject = func(c echo.Context) error {
	p := new(validators.ProjectCreate)

	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, p)
}
var CloneProject = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var DeployVercel = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var UpdateProject = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var DeleteProject = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}

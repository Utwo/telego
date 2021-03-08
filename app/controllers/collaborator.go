package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

var GetCollaboratorsByProjectId = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var AddCollaborator = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var UpdateCollaborator = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var ChangeAccess = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}
var RemoveCollaborator = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}

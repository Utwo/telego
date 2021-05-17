package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"telego/app/models"
	"telego/app/utils"
	"telego/app/validators"
)

var GetAuthenticated = func(c echo.Context) error {
	authUser := c.Get("account").(*models.Account)
	return c.JSON(http.StatusOK, authUser)
}
var LinkAnonymously = func(c echo.Context) error {
	authAnonymousId := c.Param("anonymousAuthId")
	authUser := c.Get("account").(*models.Account)
	db := c.Get("db").(*gorm.DB)

	anonymousUser := models.Account{AuthId: authAnonymousId}
	result := db.First(&anonymousUser)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) || !models.IsAnonymous(&anonymousUser) {
		return c.JSON(http.StatusBadRequest, utils.JsonError{
		Message: "Anonymous account not found",
		})
	}
	db.Model(&models.Collaborator{}).
		Where(&models.Collaborator{AccountID: &anonymousUser.Id}).
		Updates(&models.Collaborator{AccountID: &authUser.Id})
	db.Model(&models.Asset{}).
		Where(&models.Asset{UploadedByID: anonymousUser.Id}).
		Updates(&models.Asset{UploadedByID: authUser.Id})
	db.Delete(&anonymousUser)

	return c.JSON(http.StatusOK, "Hey there")
}
var UpdateAccount = func(c echo.Context) error {
	a := c.Get("body").(*validators.AccountUpdate)
	authUser := c.Get("account").(*models.Account)
	db := c.Get("db").(*gorm.DB)
	db.Model(&authUser).Updates(&models.Account{Meta: a.Meta})

	return c.JSON(http.StatusOK, authUser)
}

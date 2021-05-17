package controllers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"telego/app/limiters"
	"telego/app/models"
	"telego/app/services"
	"telego/app/utils"
	"telego/app/validators"
)

var GetCollaboratorsByProjectId = func(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	db := c.Get("db").(*gorm.DB)
	authUser := c.Get("account").(*models.Account)

	var collab models.Collaborator
	db.Model(&models.Collaborator{}).Where(&models.Collaborator{AccountID: &authUser.Id, ProjectID: id}).First(&collab)
	return c.JSON(http.StatusOK, collab)
}

var AddCollaborator = func(c echo.Context) error {
	collab := c.Get("body").(*validators.CollaboratorCreate)
	db := c.Get("db").(*gorm.DB)
	authUser := c.Get("account").(*models.Account)

	if authUser.Email == collab.Email {
		return  c.JSON(http.StatusBadRequest, utils.JsonError{
			Message: "Account cannot invite himself",
		})
	}

	var project models.Project
	result := db.First(&project, collab.ProjectId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return  echo.ErrNotFound
	}

	var inviteAccount models.Account
	newCollab := models.Collaborator{ProjectID: collab.ProjectId}
	result = db.Model(&inviteAccount).Where(&models.Account{Email: collab.Email}).First(&inviteAccount)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// We should set invitedEmail field only when inviteAccountId is not provided
		newCollab.InvitedEmail = collab.Email
	} else {
		newCollab.AccountID = &inviteAccount.Id
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := limiters.LimitCollaborators(tx, project.Id, authUser); err != nil{
			return err
		}

		if err := tx.Where(&newCollab).
			Attrs(&models.Collaborator{AccessType: collab.AccessType}).
			FirstOrCreate(&newCollab).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, utils.LimitCollabReached){
			return c.JSON(http.StatusBadRequest, utils.JsonError{Message: err.Error()})
		}
		log.Println(err)
		return echo.ErrInternalServerError
	}
	invitedByName := authUser.Name
	if invitedByName == "" {
		invitedByName = strings.Split(authUser.Email, "@")[0]
	}
	services.SendCollaborationInvitationMail(collab.Email, inviteAccount.Name, invitedByName, project.Name)

	return c.JSON(http.StatusOK, newCollab)
	//TODO: notify on websocket that a new collab was added
}

var UpdateCollaborator = func(c echo.Context) error {
	v := c.Get("body").(*validators.CollaboratorUpdate)
	db := c.Get("db").(*gorm.DB)
	authUser := c.Get("account").(*models.Account)

	collab := &models.Collaborator{Meta: v.Meta}
	db.Model(&collab).Debug().
		Where(&models.Collaborator{AccountID: &authUser.Id, ProjectID: v.ProjectId}).
		Updates(&collab)

	return c.JSON(http.StatusOK, collab)
}

var ChangeAccess = func(c echo.Context) error {
	v := c.Get("body").(*validators.CollaboratorChangeAccess)
	db := c.Get("db").(*gorm.DB)

	collab := models.Collaborator{AccessType: v.AccessType}
	if result := db.Model(&models.Collaborator{}).
		Where(&models.Collaborator{ID: v.CollaboratorId, ProjectID: v.ProjectId}).
		Updates(&collab); result.Error != nil {
		return echo.ErrInternalServerError
	}
	c.JSON(http.StatusOK, collab)

	if len(collab.InvitedEmail) > 0 {
		return nil
	}
	//TODO: notify shareDB to change the collaborator from the active websocket connections
	return nil
}

var RemoveCollaborator = func(c echo.Context) error {
	v := c.Get("body").(*validators.CollaboratorDelete)
	db := c.Get("db").(*gorm.DB)

	var collab models.Collaborator
	result := db.Model(&models.Collaborator{}).Where(&models.Collaborator{ID: v.CollaboratorId}).First(&collab)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return echo.ErrNotFound
	}

	if collab.AccessType == models.OwnerCollaborator {
		return echo.ErrForbidden
	}

	// it's important to pass both projectId and collaboratorId
	// in order to delete collaborators only from that specific project
	result = db.
		Where(&models.Collaborator{ProjectID: v.ProjectId, ID: collab.ID}).
		Delete(&models.Collaborator{})
	if result.Error != nil {
		return echo.ErrInternalServerError
	}
	if result.RowsAffected == 0 {
		return echo.ErrNotFound
	}

	c.JSON(http.StatusOK, true)

	//TODO: notify shareDB to remove the collaborator from the active websocket connections
	return nil
}

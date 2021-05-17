package middlewares

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"telego/app/models"
)

func IsProjectOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		account := c.Get("account").(*models.Account)
		db := c.Get("db").(*gorm.DB)
		projectId := c.Get("id")
		if projectId == nil {
			projectId = c.FormValue("projectId")
		}

		isProjectOwner := HaveAccess(db, projectId.(uuid.UUID), &account.Id, models.OwnerCollaborator)
		if isProjectOwner == 0 {
			return echo.ErrForbidden
		}

		return next(c)
	}
}

func CanReadProject(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		account := c.Get("account").(*models.Account)
		db := c.Get("db").(*gorm.DB)
		projectId := c.Get("id").(uuid.UUID)

		isProjectOwner := HaveAccess(db, projectId, &account.Id, models.ReadCollaborator)
		if isProjectOwner == 0 {
			return echo.ErrForbidden
		}

		return next(c)
	}
}

func CanWriteProject(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		account := c.Get("account").(*models.Account)
		db := c.Get("db").(*gorm.DB)
		projectId := c.Get("id")
		if projectId == nil {
			projectId = c.FormValue("projectId")
		}

		isProjectOwner := HaveAccess(db, projectId.(uuid.UUID), &account.Id, models.WriteCollaborator)
		if isProjectOwner == 0 {
			return echo.ErrForbidden
		}

		return next(c)
	}
}

func CanReadProjectBySlug(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		account := c.Get("account").(*models.Account)
		db := c.Get("db").(*gorm.DB)
		// is a slug
		slug := c.Param("id")

		var collaboratorCount int64
		db.Model(&models.Project{}).
			Preload("Collaborators").
			Where(&models.Project{Slug: slug, AccessType: models.PublicProject}).
			Or(db.Where(&models.Project{Slug: slug}).Where("Collaborators.accountId = ?", account.Id)).
			Count(&collaboratorCount)

		if collaboratorCount == 0 {
			return echo.ErrForbidden
		}

		return next(c)
	}
}

func HaveAccess(db *gorm.DB, projectId uuid.UUID, accountId *uuid.UUID, accessType models.AccessTypeCollaborator) int64 {
	var collaboratorCount int64
	if accessType == models.OwnerCollaborator {
		db.Model(&models.Project{}).
			Preload("Collaborators").
			Where(&models.Project{Id: projectId, AccessType: models.PublicProject}).
			Or(&models.Project{Id: projectId}, "Collaborators.accountId = ?", accountId).
			Count(&collaboratorCount)
	} else {
		accessTypeList := []models.AccessTypeCollaborator{models.OwnerCollaborator}
		if accessType == models.WriteCollaborator {
			accessTypeList = []models.AccessTypeCollaborator{models.WriteCollaborator, models.OwnerCollaborator}
		}
		db.Model(&models.Collaborator{}).
			Where(&models.Collaborator{ProjectID: projectId, AccountID: accountId}).
			Where("accessType in (?)", accessTypeList).
			Count(&collaboratorCount)
	}
	return collaboratorCount
}

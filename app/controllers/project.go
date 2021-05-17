package controllers

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"telego/app/limiters"
	"telego/app/models"
	"telego/app/services"
	"telego/app/utils"
	"telego/app/validators"
	"time"
)

var GetProjects = func(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	authAccount := c.Get("account").(*models.Account)

	// TODO: try to refactor the query below
	var p []models.Project
	db.Where("id in ((?)::uuid)", db.Table("project-collaborators").
		Select("'projectId'").
		Where("'accountId' = ?", authAccount.Id.String())).
		//TODO: Order("'createdAt' desc").
		Preload("Assets", "'assetType' = ?", models.ProjectOgImage).
		Preload("Collaborators").
		Find(&p)

	return c.JSON(http.StatusOK, p)
}
var GetBySlugOfLoggedInUser = func(c echo.Context) error {
	// https://github.com/labstack/echo/issues/1490
	// is slug not id
	slug := c.Param("id")
	db := c.Get("db").(*gorm.DB)
	var p models.Project
	result := db.Where(models.Project{Slug: slug}).
		Preload("Collaborators.Account").
		First(&p)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return echo.ErrNotFound
	}
	if result.Error != nil {
		log.Println(result.Error.Error())
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, p)
}

var CheckIfSlugIsAvailable = func(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	var projectCount int64
	db.Model(models.Project{}).
		Where(models.Project{Slug: c.Param("slug")}).
		Count(&projectCount)
	return c.JSON(http.StatusOK, projectCount)
}
var CreateProject = func(c echo.Context) error {
	p := c.Get("body").(*validators.ProjectCreate)
	db := c.Get("db").(*gorm.DB)
	authUser := c.Get("account").(*models.Account)

	accessType := models.PrivateProject
	if models.IsAnonymous(authUser) {
		accessType = models.PublicProject
	}

	timestamp := time.Now().Unix()
	project := models.Project{
		Name:       p.Name,
		Slug:       slug.Make(p.Name) + "-" + strconv.Itoa(int(timestamp)),
		AccessType: accessType,
		Collaborators: []models.Collaborator{{
			Account:    authUser,
			AccessType: models.OwnerCollaborator,
		}},
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := limiters.LimitProjects(tx, authUser); err != nil {
			return err
		}
		if err := tx.Create(&project).Preload("Collaborators").Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, utils.LimitProjectReached){
			return c.JSON(http.StatusBadRequest, utils.JsonError{Message: err.Error()})
		}
		log.Println(err)
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, project)
}
var CloneProject = func(c echo.Context) error {
	id := c.Get("id").(uuid.UUID)
	db := c.Get("db").(*gorm.DB)
	//gs := c.Get("gs").(*services.GoogleStorage)
	authUser := c.Get("account").(*models.Account)

	project := models.Project{Id: id}
	result := db.First(&project)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return echo.ErrNotFound
	}

	timestamp := time.Now().Unix()
	slug := slug.Make(project.Name + "-" + strconv.Itoa(int(timestamp)))
	name := "Clone of " + project.Name
	newProject := models.Project{
		Name:       name,
		Slug:       slug,
		AccessType: project.AccessType,
		Collaborators: []models.Collaborator{{
			Account:    authUser,
			AccessType: models.OwnerCollaborator,
		}}}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := limiters.LimitProjects(tx, authUser); err != nil {
			return err
		}
		if err := tx.Create(&newProject).Preload("Collaborators").Error; err != nil {
			return err
		}
		if err := tx.Exec(
			"INSERT INTO snapshots(doc_id, collection, doc_type, version, data, \"updatedAt\") SELECT @destProjectId AS doc_id, collection, doc_type, version, data, NOW() FROM snapshots WHERE doc_id = @srcProjectId and collection = 'project';",
			sql.Named("srcProjectId", project.Id),
			sql.Named("destProjectId", newProject.Id),
		).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, utils.LimitProjectReached){
			return c.JSON(http.StatusBadRequest, utils.JsonError{Message: err.Error()})
		}
		log.Println(err)
		return echo.ErrInternalServerError
	}
	//TODO: clone assets
	//services.GoogleStorage.CopyFiles(gs, project.Id, '', newProject.Id)
	return c.JSON(http.StatusOK, &newProject)
}
var UpdateProject = func(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	id := c.Get("id").(uuid.UUID)
	p := c.Get("body").(*validators.ProjectUpdate)
	authUser := c.Get("account").(*models.Account)

	if len(p.Slug) > 0 {
		var slugCount int64
		db.Model(&models.Project{}).Where(models.Project{Slug: p.Slug}).Count(&slugCount)
		if slugCount > 0 {
			return c.JSON(http.StatusBadRequest, utils.JsonError{Message: "Slug is already taken"})
		}
	}

	if len(p.AccessType) > 0 && models.IsAnonymous(authUser) {
		return echo.ErrBadRequest
	}
	project := &models.Project{
		Id:         id,
		Name:       p.Name,
		Slug:       p.Slug,
		AccessType: p.AccessType,
	}
	if result := db.Model(&project).Updates(&project); result.Error != nil {
		log.Println(result.Error.Error())
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, project)
}
var DeleteProject = func(c echo.Context) error {
	id := c.Param("id")
	db := c.Get("db").(*gorm.DB)
	if result := db.Delete(&models.Project{}, &id); result.Error != nil {
		log.Println(result.Error)
		return echo.ErrInternalServerError
	}
	_ = c.JSON(http.StatusOK, true)

	// TODO:
	//VercelDeploy.deleteProject(project.slug)
	gs := c.Get("gs").(*services.GoogleStorage)
	ctx := c.Request().Context()
	if err := gs.DeleteFolders(ctx, id); err != nil {
		log.Println(err)
	}
	// notify shareDB to remove all active websocket connections from that project
	//publishToRedis(REDIS_CHANNELS.deleteProject, { projectId: id })

	return nil
}
var DeployVercel = func(c echo.Context) error {
	return c.String(http.StatusOK, "Hey there")
}

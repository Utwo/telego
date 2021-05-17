package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"net/http"
	"telego/app/config"
	"telego/app/limiters"
	"telego/app/models"
	"telego/app/services"
	"telego/app/utils"
	"telego/app/validators"
)

var UploadAsset = func(c echo.Context) error {
	gs := c.Get("gs").(*services.GoogleStorage)
	db := c.Get("db").(*gorm.DB)
	authUser := c.Get("account").(*models.Account)
	v := c.Get("body").(*validators.AssetCreate)

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	maxUploadBatch := models.GetLimitations(authUser).MaxUploadBatch
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, utils.JsonError{Message: "No file uploaded"})
	}
	if len(files) > int(maxUploadBatch) {
		return c.JSON(http.StatusBadRequest, utils.JsonError{Message: fmt.Sprintf("Max upload batch reached: %v", maxUploadBatch)})
	}
	var uploadResponse []models.Asset
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		objectName, err := gs.UploadFile(c.Request().Context(), src, v.ProjectId.String())
		if err != nil {
			log.Println(err)
			rollbackUploads(c.Request().Context(), gs, uploadResponse, v.ProjectId)
			return echo.ErrInternalServerError
		}
		//TODO: mimeType, err := mimetype.DetectReader(src)
		//if err != nil {
		//	log.Println(err)
		//	rollbackUploads(c.Request().Context(), gs, uploadResponse, v.ProjectId)
		//	return echo.ErrInternalServerError
		//}

		uploadResponse = append(uploadResponse, models.Asset{
			Path:         *objectName,
			ProjectID:    v.ProjectId,
			MimeType:     "",
			Size:         file.Size,
			Bucket:       config.Config.GCloud.GcloudStorageBucket,
			AssetType:    v.AssetType,
			UploadedByID: authUser.Id,
		})
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := limiters.LimitAssets(tx, uploadResponse, authUser); err != nil {
			return err
		}
		if err := tx.Create(&uploadResponse).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		rollbackUploads(c.Request().Context(), gs, uploadResponse, v.ProjectId)
		if errors.Is(err, utils.LimitAssetReached) || errors.Is(err, utils.SizeAssetReached) {
			return c.JSON(http.StatusBadRequest, err)
		}
		log.Println(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, uploadResponse)
}

var DuplicateAsset = func(c echo.Context) error {
	gs := c.Get("gs").(*services.GoogleStorage)
	db := c.Get("db").(*gorm.DB)
	authUser := c.Get("account").(*models.Account)
	v := c.Get("body").(*validators.AssetDuplicate)

	var assets []models.Asset
	if err := db.Where(models.Asset{ProjectID: v.ProjectId}).Where("path in (?)", v.Paths).Find(&assets).Error; err != nil {
		log.Println(err)
		return echo.ErrInternalServerError
	}

	var newAssets []models.Asset
	var newPaths []string
	var oldPaths []string
	for _, asset := range assets {
		newPath := uuid.New()
		oldPaths = append(oldPaths, asset.Path.String())
		newPaths = append(newPaths, newPath.String())
		newAssets = append(newAssets, models.Asset{
			Path:         newPath,
			ProjectID:    v.ProjectId,
			UploadedByID: authUser.Id,
			MimeType:     asset.MimeType,
			Size:         asset.Size,
			Bucket:       asset.Bucket,
			AssetType:    asset.AssetType,
		})
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newAssets).Error; err != nil {
			return err
		}
		if _, err := gs.CopyFiles(c.Request().Context(), v.ProjectId.String(), oldPaths, v.ProjectId.String(), newPaths); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Println(err)
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, newAssets)
}

var DeleteAsset = func(c echo.Context) error {
	gs := c.Get("gs").(*services.GoogleStorage)
	db := c.Get("db").(*gorm.DB)
	v := c.Get("body").(*validators.AssetDelete)

	if err := db.Where(models.Asset{ProjectID: v.ProjectId}).Where("path in (?)", v.Paths).Delete(&models.Asset{}).Error; err != nil {
		log.Println(err)
		return echo.ErrInternalServerError
	}
	var paths []string
	for _, path := range v.Paths {
		paths = append(paths, v.ProjectId.String()+"/"+path)
	}
	if err := gs.DeleteFiles(c.Request().Context(), paths); err != nil {
		log.Println(err)
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, true)
}

func rollbackUploads(ctx context.Context, gs *services.GoogleStorage, uploadResponse []models.Asset, projectId uuid.UUID) {
	var assetsPath []string
	for _, asset := range uploadResponse {
		assetsPath = append(assetsPath, projectId.String()+"/"+asset.Path.String())
	}
	if err := gs.DeleteFiles(ctx, assetsPath); err != nil {
		log.Println(err)
	}
}

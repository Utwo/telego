package app

import (
	"github.com/labstack/echo/v4"
	"telego/app/controllers"
	"telego/app/middlewares"
)

var InitRoutes = func(e *echo.Echo) {
	accountGroup := e.Group("/account")
	accountGroup.Use(middlewares.Auth(false))
	accountGroup.GET("", controllers.GetAuthenticated)
	accountGroup.PUT("", controllers.UpdateAccount)
	accountGroup.POST("/link-anonymous/:anonymousAuthId", controllers.LinkAnonymously)

	projectGroup := e.Group("/project")
	projectGroup.GET("/all", controllers.GetProjects)
	projectGroup.GET("/:slug", controllers.GetBySlugOfLoggedInUser)
	projectGroup.GET("/:id/stats", controllers.GetAssetStats)
	projectGroup.GET("/slug-availability/:slug", controllers.CheckIfSlugIsAvailable)
	projectGroup.POST("", controllers.CreateProject)
	projectGroup.POST("/:id/clone", controllers.CloneProject)
	projectGroup.POST("/:id/deploy", controllers.DeployVercel)
	projectGroup.PUT("/:id", controllers.UpdateProject)
	projectGroup.DELETE("/:id", controllers.DeleteProject)

	collaboratorGroup := e.Group("/collaborator")
	collaboratorGroup.GET("/project/:id", controllers.GetCollaboratorsByProjectId)
	collaboratorGroup.POST("", controllers.AddCollaborator)
	collaboratorGroup.PUT("", controllers.UpdateCollaborator)
	collaboratorGroup.PUT("/change-access", controllers.ChangeAccess)
	collaboratorGroup.DELETE("", controllers.RemoveCollaborator)

	assetGroup := e.Group("/asset")
	assetGroup.POST("", controllers.UploadAsset)
	assetGroup.POST("/duplicate", controllers.DuplicateAsset)
	assetGroup.DELETE("", controllers.DeleteAsset)
}

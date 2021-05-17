package app

import (
	"github.com/labstack/echo/v4"
	"telego/app/controllers"
	"telego/app/middlewares"
	"telego/app/validators"
)

var InitRoutes = func(e *echo.Echo) {
	accountGroup := e.Group("/account")
	accountGroup.Use(middlewares.Auth(false))
	accountGroup.GET("", controllers.GetAuthenticated)

	accountGroup.PUT("",
		controllers.UpdateAccount,
		middlewares.ValidateRequest(new(validators.AccountUpdate)),
	)
	accountGroup.POST("/link-anonymous/:anonymousAuthId", controllers.LinkAnonymously)

	projectGroup := e.Group("/project")
	projectGroup.GET("/all",
		controllers.GetProjects,
		middlewares.Auth(false),
	)
	projectGroup.GET("/:id",
		controllers.GetBySlugOfLoggedInUser,
		middlewares.Auth(true),
		middlewares.CanReadProjectBySlug,
	)
	projectGroup.GET("/slug-availability/:slug",
		controllers.CheckIfSlugIsAvailable,
		middlewares.Auth(false),
	)
	projectGroup.POST("",
		controllers.CreateProject,
		middlewares.Auth(false),
		middlewares.ValidateRequest(new(validators.ProjectCreate)),
	)
	projectGroup.POST("/:id/clone",
		controllers.CloneProject,
		middlewares.ParseId,
		middlewares.Auth(false),
		middlewares.CanReadProject,
	)
	projectGroup.POST("/:id/deploy",
		controllers.DeployVercel,
		middlewares.ParseId,
		middlewares.Auth(false),
		middlewares.CanWriteProject,
	)
	projectGroup.PUT("/:id",
		controllers.UpdateProject,
		middlewares.ParseId,
		middlewares.Auth(false),
		middlewares.ValidateRequest(new(validators.ProjectUpdate)),
		middlewares.CanWriteProject,
	)
	projectGroup.DELETE("/:id",
		controllers.DeleteProject,
		middlewares.ParseId,
		middlewares.Auth(false),
		middlewares.IsProjectOwner,
	)

	collaboratorGroup := e.Group("/collaborator")
	collaboratorGroup.Use(middlewares.Auth(false))
	collaboratorGroup.GET("/project/:id",
		controllers.GetCollaboratorsByProjectId,
		middlewares.ParseId,
		middlewares.CanReadProject,
	)
	collaboratorGroup.POST("",
		controllers.AddCollaborator,
		middlewares.CanWriteProject,
		middlewares.ValidateRequest(new(validators.CollaboratorCreate)),
	)
	collaboratorGroup.PUT("",
		controllers.UpdateCollaborator,
		middlewares.CanWriteProject,
		middlewares.ValidateRequest(new(validators.CollaboratorUpdate)),
	)
	collaboratorGroup.PUT("/change-access",
		controllers.ChangeAccess,
		middlewares.IsProjectOwner,
		middlewares.ValidateRequest(new(validators.CollaboratorChangeAccess)),
	)
	collaboratorGroup.DELETE("",
		controllers.RemoveCollaborator,
		middlewares.CanWriteProject,
		middlewares.ValidateRequest(new(validators.CollaboratorDelete)),
	)

	assetGroup := e.Group("/asset")
	assetGroup.Use(middlewares.Auth(false))
	assetGroup.POST("",
		controllers.UploadAsset,
		middlewares.CanWriteProject,
		middlewares.ValidateRequest(new(validators.AssetCreate)),
	)
	assetGroup.POST("/duplicate",
		controllers.DuplicateAsset,
		middlewares.CanWriteProject,
		middlewares.ValidateRequest(new(validators.AssetDuplicate)),
	)
	assetGroup.DELETE("",
		controllers.DeleteAsset,
		middlewares.CanWriteProject,
		middlewares.ValidateRequest(new(validators.AssetDelete)),
	)
}

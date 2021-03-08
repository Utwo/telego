package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"telego/app"
	"telego/app/config"
	"telego/app/models"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return echo.NewHTTPError(http.StatusInternalServerError, cv.validator.Struct(i).Error())
}

func main() {
	dsn := "host=" + config.Config.Postgres.PostgresServer + " user=" + config.Config.Postgres.PostgresUser + " password=" + config.Config.Postgres.PostgresPassword + " dbname=" + config.Config.Postgres.PostgresDb + " port=" + config.Config.Postgres.PostgresPortExternal
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(models.Account{}, models.Project{}, models.Collaborator{}, models.Asset{})
	if err != nil {
		panic("failed to run migrations")
	}

	ctx := context.Background()
	firebaseApp, err := firebase.NewApp(ctx, nil)
	if err != nil {
		panic(err)
	}

	client, err := firebaseApp.Auth(ctx)
	if err != nil {
		panic("failed to use firebase Auth")
	}

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	//e.Use(middleware.Logger())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("firebase", client)
			c.Set("firebaseCtx", ctx)
			c.Set("db", db)
			return next(c)
		}
	})
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))

	app.InitRoutes(e)

	e.Logger.Fatal(e.Start(":" + config.Config.Port))
}

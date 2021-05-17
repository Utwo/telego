package main

import (
	"context"
	"database/sql"
	"errors"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/iancoleman/strcase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"telego/app"
	"telego/app/config"
	"telego/app/middlewares"
	"telego/app/services"
	"telego/app/validators"
)

type customNameStrategy struct{ schema.NamingStrategy }

func (ns customNameStrategy) ColumnName(table, column string) string {
	return strcase.ToLowerCamel(column)
}

func main() {
	dsn := "postgres://" + config.Config.Postgres.PostgresUser + ":" + config.Config.Postgres.PostgresPassword + "@" + config.Config.Postgres.PostgresServer + ":" + config.Config.Postgres.PostgresPortExternal + "/" + config.Config.Postgres.PostgresDb + "?sslmode=disable"
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("cannot connect to the database: %v", err))
	}
	driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://app/db/migrations",
		"postgres", driver)
	if err != nil {
		panic(fmt.Errorf("cannot create new instance for migrations: %v", err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(fmt.Errorf("cannot run migrations: %v", err))
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{NamingStrategy: customNameStrategy{schema.NamingStrategy{}}})
	if err != nil {
		panic("failed to connect database")
	}

	ctx := context.Background()
	firebaseApp, err := firebase.NewApp(ctx, nil)
	if err != nil {
		panic(fmt.Errorf("failed to use firbase App: %v", err))
	}

	client, err := firebaseApp.Auth(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to use firebase Auth: %v", err))
	}

	gs, err := services.NewGoogleStorage(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to use google storage bucket: %v", err))
	}

	e := echo.New()
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))

	customValidator := validator.New()
	customValidator.RegisterValidation("slug", validators.SlugValidator)
	e.Validator = &middlewares.CustomValidator{Validator: customValidator}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("firebase", client)
			c.Set("firebaseCtx", ctx)
			c.Set("db", db)
			c.Set("gs", gs)
			return next(c)
		}
	})

	app.InitRoutes(e)

	e.Logger.Fatal(e.Start(":" + config.Config.Port))
}

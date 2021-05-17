package config

import (
	"os"
)

type postgres struct {
	PostgresServer            string
	PostgresPortExternal      string
	PostgresUser              string
	PostgresPassword          string
	PostgresDb                string
	PostgresConnectionMaxPool string
}

type email struct {
	SendgridApiKey string
	TeleportEmail  string
}

type gcloud struct {
	GoogleCloudProject  string
	GcloudStorageBucket string
}

var Config = struct {
	AppEnv   string
	Port     string
	Email    email
	GCloud   gcloud
	Postgres postgres
}{
	AppEnv: os.Getenv("APP_ENV"),
	Port:   os.Getenv("PORT"),
	Email: email{
		SendgridApiKey: os.Getenv("SENDGRID_API_KEY"),
		TeleportEmail:  "hello@teleporthq.io",
	},
	GCloud: gcloud{
		GoogleCloudProject:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
		GcloudStorageBucket: os.Getenv("GCLOUD_STORAGE_BUCKET"),
	},
	Postgres: postgres{
		PostgresServer:            os.Getenv("POSTGRES_SERVER"),
		PostgresPortExternal:      os.Getenv("POSTGRES_PORT_EXTERNAL"),
		PostgresUser:              os.Getenv("POSTGRES_USER"),
		PostgresPassword:          os.Getenv("POSTGRES_PASSWORD"),
		PostgresDb:                os.Getenv("POSTGRES_DB"),
		PostgresConnectionMaxPool: os.Getenv("POSTGRES_CONNECTION_MAX_POOL"),
	},
}

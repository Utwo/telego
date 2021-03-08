package models

import (
	"database/sql/driver"
	"github.com/google/uuid"
	"time"
)

type Asset struct {
	Path         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID    uuid.UUID `gorm:"primaryKey"`
	Project      Project   `gorm:"constraint:OnDelete:CASCADE;"`
	MimeType     string    `gorm:"type:string;not null;size:50"`
	Size         uint16    `gorm:"not null"`
	Bucket       string    `gorm:"type:string;not null;size:40"`
	AssetType    assetType `gorm:"not null" ;sql:"type:asset_type"`
	UploadedByID uuid.UUID
	UploadedBy   Account `gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type assetType string

const (
	image          assetType = "image"
	favIcon        assetType = "favIcon"
	pageOgImage    assetType = "pageOgImage"
	projectOgImage assetType = "projectOgImage"
)

func (a *assetType) Scan(value interface{}) error {
	*a = assetType(value.([]byte))
	return nil
}

func (a assetType) Value() (driver.Value, error) {
	return string(a), nil
}

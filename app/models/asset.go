package models

import (
	"database/sql/driver"
	"github.com/google/uuid"
	"time"
)

type Asset struct {
	Path         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"path,omitempty"`
	ProjectID    uuid.UUID `gorm:"primaryKey" json:"project_id,omitempty"`
	Project      Project   `gorm:"constraint:OnDelete:CASCADE;" json:"project,omitempty"`
	MimeType     string    `gorm:"type:string;not null;size:50" json:"mime_type,omitempty"`
	Size         int64    `gorm:"not null" json:"size,omitempty"`
	Bucket       string    `gorm:"type:string;not null;size:40" json:"bucket,omitempty"`
	AssetType    AssetType `gorm:"not null" ;sql:"type:asset_type" json:"asset_type,omitempty"`
	UploadedByID uuid.UUID `json:"uploaded_by_id,omitempty"`
	UploadedBy   Account   `gorm:"constraint:OnDelete:CASCADE;" json:"uploaded_by,omitempty"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type AssetType string

const (
	Image          AssetType = "image"
	FavIcon        AssetType = "favIcon"
	PageOgImage    AssetType = "pageOgImage"
	ProjectOgImage AssetType = "projectOgImage"
)

func (a *AssetType) Scan(value interface{}) error {
	*a = AssetType(value.(string))
	return nil
}

func (a AssetType) Value() (driver.Value, error) {
	return string(a), nil
}

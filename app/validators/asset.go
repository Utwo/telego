package validators

import (
	"github.com/google/uuid"
	"telego/app/models"
)

type AssetCreate struct {
	AssetType models.AssetType `form:"assetType" validate:"oneof=image fav-icon page-og-image project-og-image"`
	ProjectId uuid.UUID `form:"projectId" validate:"required"` //TODO: uuid4
}

type AssetDuplicate struct {
	ProjectId uuid.UUID `json:"projectId" validate:"required,uuid4"` //TODO: uuid4
	Paths []string `json:"paths" validate:"required,dive,uuid4,required"` //TODO: uuid4
}

type AssetDelete struct {
	ProjectId uuid.UUID `json:"projectId" validate:"required"` //TODO: uuid4
	Paths []string `json:"paths" validate:"required"` //TODO: uuid4
}




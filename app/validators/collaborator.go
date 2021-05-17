package validators

import (
	"github.com/google/uuid"
	"telego/app/models"
)

type CollaboratorCreate struct {
	Email      string                        `json:"email" validate:"required,email"`
	ProjectId  uuid.UUID                     `json:"projectId" validate:"required"` //TODO uuid4
	AccessType models.AccessTypeCollaborator `json:"accessType" validate:"required,oneof=read write"`
}

type CollaboratorChangeAccess struct {
	CollaboratorId uuid.UUID                     `json:"collaboratorId" validate:"required"` //TODO uuid4
	ProjectId      uuid.UUID                     `json:"projectId" validate:"required"`      //TODO uuid4
	AccessType     models.AccessTypeCollaborator `json:"accessType" validate:"required,oneof=read write"`
}

type CollaboratorUpdate struct {
	Meta      map[string]interface{} `json:"meta" validate:"required"`
	ProjectId uuid.UUID              `json:"projectId" validate:"required"` //TODO uuid4
}

type CollaboratorDelete struct {
	CollaboratorId uuid.UUID `json:"collaboratorId" validate:"required"` //TODO uuid4
	ProjectId      uuid.UUID `json:"projectId" validate:"required"`      //TODO uuid4
}

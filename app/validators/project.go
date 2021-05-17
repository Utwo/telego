package validators

import "telego/app/models"

type ProjectCreate struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type ProjectUpdate struct {
	Name string `json:"name" validate:"omitempty,min=3,max=100"`
	Slug string `json:"slug" validate:"omitempty,slug,min=3,max=20"`
	AccessType models.AccessTypeProject `json:"accessType" validate:"omitempty,oneof=public private"`
}

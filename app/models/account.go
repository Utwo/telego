package models

import (
	"github.com/google/uuid"
	"telego/app/constants"
	"telego/app/models/utils"
	"time"
)

type Account struct {
	Id                    uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id,omitempty"`
	AuthId                string         `gorm:"not null;size:100" json:"auth_id,omitempty"`
	Name                  string         `gorm:"size:255" json:"name,omitempty"`
	Email                 string         `gorm:"not null;size:100;unique" json:"email,omitempty"`
	Picture               string         `gorm:"not null" json:"picture,omitempty"`
	Meta                  utils.JSONB    `gorm:"type:jsonb;default:(-);" json:"meta"`
	IsAnonymous           bool           `gorm:"->;type:GENERATED ALWAYS AS (email is NULL);default:(-);" json:"is_anonymous,omitempty"`
	Collaborators         []Collaborator `gorm:"constraint:OnDelete:CASCADE;" json:"collaborators,omitempty"`
	CollaboratingProjects []Project      `gorm:"many2many:collaborators;" json:"collaborating_projects,omitempty"`
	CreatedAt             time.Time      `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt             time.Time      `gorm:"not null" json:"updated_at,omitempty"`
}

func IsAnonymous(account *Account) bool {
	return account.Email == ""
}

func GetLimitations(account *Account) constants.PlanLimitations {
	if account.IsAnonymous {
		return constants.AnonymousPlanLimitations
	}
	return constants.DefaultPlanLimitations
}

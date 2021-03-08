package models

import (
	"github.com/google/uuid"
	"telego/app/models/utils"
	"time"
)

type Account struct {
	Id                    uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthId                string         `gorm:"not null;size:100"`
	Name                  string         `gorm:"size:255"`
	Email                 string         `gorm:"not null;size:100;unique"`
	Picture               string         `gorm:"not null"`
	Meta                  utils.JSONB    `gorm:"type:jsonb"`
	IsAnonymous           bool           `gorm:"default:false;"`
	Collaborators         []Collaborator `gorm:"constraint:OnDelete:CASCADE;"`
	CollaboratingProjects []Project      `gorm:"many2many:collaborators;"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

package models

import (
	"database/sql/driver"
	"github.com/google/uuid"
	"time"
)

type Project struct {
	Id                    uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid()" json:"id,omitempty"`
	Name                  string            `gorm:"not null;size:100" json:"name,omitempty"`
	Slug                  string            `gorm:"not null;unique;size:100" json:"slug,omitempty"`
	AccessType            AccessTypeProject `gorm:"not null" ;sql:"type:access_type" json:"access_type,omitempty"`
	Collaborators         []Collaborator    `gorm:"constraint:OnDelete:CASCADE;" json:"collaborators,omitempty"`
	CollaboratorsAccounts []Account         `gorm:"many2many:collaborators;" json:"collaborators_accounts,omitempty"`
	Assets                []Asset           `json:"assets,omitempty"`
	CreatedAt             time.Time         `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt             time.Time         `gorm:"not null" json:"updated_at,omitempty"`
}

type AccessTypeProject string

const (
	PublicProject  AccessTypeProject = "public"
	PrivateProject AccessTypeProject = "private"
)

func (a *AccessTypeProject) Scan(value interface{}) error {
	*a = AccessTypeProject(value.(string))
	return nil
}

func (a AccessTypeProject) Value() (driver.Value, error) {
	return string(a), nil
}
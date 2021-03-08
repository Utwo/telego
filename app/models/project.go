package models

import (
	"database/sql/driver"
	"github.com/google/uuid"
	"time"
)

type Project struct {
	Id                    uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid()"`
	Name                  string            `gorm:"not null;size:100"`
	Slug                  string            `gorm:"not null;unique;size:100"`
	AccessType            accessTypeProject `gorm:"not null" ;sql:"type:access_type"`
	Collaborators         []Collaborator    `gorm:"constraint:OnDelete:CASCADE;"`
	CollaboratorsAccounts []Account         `gorm:"many2many:collaborators;"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type accessTypeProject string

const (
	public  accessTypeProject = "public"
	private accessTypeProject = "private"
)

func (a *accessTypeProject) Scan(value interface{}) error {
	*a = accessTypeProject(value.([]byte))
	return nil
}

func (a accessTypeProject) Value() (driver.Value, error) {
	return string(a), nil
}

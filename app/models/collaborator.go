package models

import (
	"database/sql/driver"
	"github.com/google/uuid"
	"telego/app/models/utils"
	"time"
)

type Collaborator struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AccountID    uuid.UUID
	Account      Account
	ProjectID    uuid.UUID
	Project      Project
	InvitedEmail string                 `gorm:"size:100"`
	Meta         utils.JSONB            `gorm:"type:jsonb"`
	AccessType   accessTypeCollaborator `gorm:"not null" ;sql:"type:access_type"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type accessTypeCollaborator string

const (
	owner accessTypeCollaborator = "owner"
	write accessTypeCollaborator = "write"
	read  accessTypeCollaborator = "read"
)

func (a *accessTypeCollaborator) Scan(value interface{}) error {
	*a = accessTypeCollaborator(value.([]byte))
	return nil
}

func (a accessTypeCollaborator) Value() (driver.Value, error) {
	return string(a), nil
}

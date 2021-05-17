package models

import (
	"database/sql/driver"
	"github.com/google/uuid"
	"telego/app/models/utils"
	"time"
)

type Collaborator struct {
	ID           uuid.UUID              `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id,omitempty"`
	AccountID    *uuid.UUID             `gorm:"column:accountId" json:"account_id,omitempty"`
	Account      *Account               `json:"account,omitempty"`
	ProjectID    uuid.UUID              `gorm:"column:projectId;not null" json:"project_id,omitempty"`
	Project      Project                `json:"project,omitempty"`
	InvitedEmail string                 `gorm:"size:100" json:"invited_email,omitempty"`
	Meta         utils.JSONB            `gorm:"type:jsonb" json:"meta,omitempty"`
	AccessType   AccessTypeCollaborator `gorm:"not null" ;sql:"type:access_type" json:"access_type,omitempty"`
	CreatedAt    time.Time              `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt    time.Time              `gorm:"not null" json:"updated_at,omitempty"`
}

func (Collaborator) TableName() string {
	return "project-collaborators"
}

type AccessTypeCollaborator string

const (
	OwnerCollaborator AccessTypeCollaborator = "owner"
	WriteCollaborator AccessTypeCollaborator = "write"
	ReadCollaborator  AccessTypeCollaborator = "read"
)

func (a *AccessTypeCollaborator) Scan(value interface{}) error {
	*a = AccessTypeCollaborator(value.(string))
	return nil
}

func (a AccessTypeCollaborator) Value() (driver.Value, error) {
	return string(a), nil
}

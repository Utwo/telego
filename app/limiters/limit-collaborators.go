package limiters

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"telego/app/models"
	"telego/app/utils"
)

func LimitCollaborators(tx *gorm.DB, projectId uuid.UUID, account *models.Account) error {
	var projectCollabCount int64
	tx.Model(&models.Collaborator{}).
		Where(models.Collaborator{ProjectID: projectId}).
		Count(&projectCollabCount)

	if projectCollabCount >= int64(models.GetLimitations(account).MaxCollaboratorsInProject) {
		return utils.LimitCollabReached
	}
	return nil
}

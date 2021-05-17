package limiters

import (
	"gorm.io/gorm"
	"telego/app/models"
	"telego/app/utils"
)

func LimitProjects(tx *gorm.DB, account *models.Account) error {
	var projectCount int64
	tx.Model(&models.Collaborator{}).
		Where(models.Collaborator{AccountID: &account.Id, AccessType: models.OwnerCollaborator}).
		Count(&projectCount)

	if uint8(projectCount) >= models.GetLimitations(account).MaxProjectsForUser {
		return utils.LimitProjectReached
	}
	return nil
}

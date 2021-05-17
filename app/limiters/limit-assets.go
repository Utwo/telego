package limiters

import (
	"gorm.io/gorm"
	"telego/app/models"
	"telego/app/utils"
)

func LimitAssets(tx *gorm.DB, assets []models.Asset, account *models.Account) error {
	if len(assets) == 0 || (len(assets) == 1 && assets[0].AssetType == models.FavIcon) {
		return nil
	}
	projectId := assets[0].ProjectID
	var assetCount int64
	tx.Model(&models.Asset{}).
		Where(models.Asset{ProjectID: projectId}).
		Count(&assetCount)

	if assetCount+int64(len(assets)) >= int64(models.GetLimitations(account).MaxAssetsForProject) {
		return utils.LimitAssetReached
	}

	var currentSize int64
	_ = tx.Table("assets").
		Select("sum(size)").
		Row().
		Scan(&currentSize)

	var allAssetSize int64
	for _, asset := range assets {
		allAssetSize += asset.Size
	}
	if currentSize+allAssetSize > models.GetLimitations(account).MaxAssetsSizeForProject {
		return utils.SizeAssetReached
	}
	return nil
}

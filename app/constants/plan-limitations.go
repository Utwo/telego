package constants

type planLimitations struct {
	MaxCollaboratorsInProject         uint8
	MaxRealTimeCollaboratorsInProject uint8
	MaxProjectsForUser                uint8
	MaxAssetsSizeForProject           uint16
	MaxAssetsForProject               uint8
	MaxUploadBatch                    uint8
}

var DefaultPlanLimitations = planLimitations{
	MaxCollaboratorsInProject:         20,
	MaxRealTimeCollaboratorsInProject: 5,
	MaxProjectsForUser:                20,
	MaxAssetsSizeForProject:           200 * 1024 * 1024,
	MaxAssetsForProject:               400,
	MaxUploadBatch:                    20,
}

var AnonymousPlanLimitations = planLimitations{
	MaxCollaboratorsInProject:         0,
	MaxRealTimeCollaboratorsInProject: 0,
	MaxProjectsForUser:                3,
	MaxAssetsSizeForProject:           0,
	MaxAssetsForProject:               0,
	MaxUploadBatch:                    0,
}

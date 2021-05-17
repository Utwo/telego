package constants

type PlanLimitations struct {
	MaxCollaboratorsInProject         uint8
	MaxRealTimeCollaboratorsInProject uint8
	MaxProjectsForUser                uint8
	MaxAssetsSizeForProject           int64
	MaxAssetsForProject               uint16
	MaxUploadBatch                    uint8
}

var DefaultPlanLimitations = PlanLimitations{
	MaxCollaboratorsInProject:         20,
	MaxRealTimeCollaboratorsInProject: 5,
	MaxProjectsForUser:                20,
	MaxAssetsSizeForProject:           200 * 1024 * 1024,
	MaxAssetsForProject:               400,
	MaxUploadBatch:                    20,
}

var AnonymousPlanLimitations = PlanLimitations{
	MaxCollaboratorsInProject:         0,
	MaxRealTimeCollaboratorsInProject: 0,
	MaxProjectsForUser:                3,
	MaxAssetsSizeForProject:           0,
	MaxAssetsForProject:               0,
	MaxUploadBatch:                    0,
}

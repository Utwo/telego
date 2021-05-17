package validators

import "telego/app/models/utils"

type AccountUpdate struct {
	Meta utils.JSONB `json:"meta" validate:"required,structonly"`
}

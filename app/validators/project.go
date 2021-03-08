package validators

type ProjectCreate struct {
	Name string `json:"name" validate:"required"`
}

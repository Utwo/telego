package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
)

func SlugValidator(fl validator.FieldLevel) bool {

	if slug.IsSlug(fl.Field().String()) {
		return true
	}

	return false
}


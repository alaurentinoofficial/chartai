package databases_validations

import (
	"regexp"

	core_validations "github.com/alaurentinoofficial/chartai/internal/core/validations"
	"github.com/go-playground/validator/v10"
)

func RegisterValidations(
	structValidator *core_validations.StructValidator,
) {
	var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9-_\s]+$`)
	structValidator.Register("DatabaseName", "Invalid Database Name", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return nameRegex.MatchString(value)
	})
}

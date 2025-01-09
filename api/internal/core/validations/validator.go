package core_validations

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	validate *validator.Validate
}

func NewStructValidator() *StructValidator {
	return &StructValidator{
		validate: validator.New(),
	}
}

func (v *StructValidator) Register(name string, errorMessage string, validation validator.Func) {
	errorMessages[name] = errorMessage
	v.validate.RegisterValidation(name, validation)
}

func (v *StructValidator) Validate(value any) *FormErrors {
	if errs := v.validate.Struct(value); errs != nil {
		return ConvertToValidationFormErrors(errs.(validator.ValidationErrors))
	}

	return nil
}

func CleanNamespace(namespace string) string {
	result := strings.SplitN(namespace, ".", 2)
	return result[len(result)-1]
}

func ConvertToValidationFormErrors(validationErrors validator.ValidationErrors) *FormErrors {
	var response FormErrors

	for _, errorField := range validationErrors {
		response = append(response,
			NewFieldError(
				CleanNamespace(errorField.Namespace()),
				errorField.Tag(),
				errorField.Value()),
		)
	}

	return &response
}

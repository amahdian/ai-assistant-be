package binding

import (
	"github.com/go-playground/validator/v10"
)

func registerCustomBinding(v *validator.Validate) {
	//registerCustomFieldValidator(v, &ExistsBinding{mStg: mStg})
}

func registerCustomFieldValidator(v *validator.Validate, fieldValidator CustomFieldBinding) {
	_ = v.RegisterValidation(fieldValidator.Tag(), fieldValidator.Validate)
	registerFieldTranslator(fieldValidator.Tag(), fieldValidator.Translate)
}

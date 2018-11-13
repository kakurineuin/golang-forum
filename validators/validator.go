package validators

import (
	valid "github.com/asaskevich/govalidator"
)

type customValidator struct {
}

func (cv *customValidator) Validate(i interface{}) error {
	_, err := valid.ValidateStruct(i)
	return err
}

func InitValidator() customValidator {
	return customValidator{}
}

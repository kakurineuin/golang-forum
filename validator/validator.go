package validator

import (
	valid "github.com/asaskevich/govalidator"
)

// CustomValidator 自訂驗證器。
type CustomValidator struct {
}

// Validate 驗證 struct。
func (cv *CustomValidator) Validate(i interface{}) error {
	_, err := valid.ValidateStruct(i)
	return err
}

// InitValidator 初始化驗證器。
func InitValidator() CustomValidator {
	return CustomValidator{}
}

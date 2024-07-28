package validator

import (
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// mobile 验证手机号码
func registerMobile(validate *validator.Validate, trans ut.Translator) {
	if err := validate.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
		ok, _ := regexp.MatchString(`^(13|14|15|16|17|18|19)[0-9]{9}$`, fl.Field().String())
		return ok
	}); err != nil {
		panic(err)
	}
	if err := validate.RegisterTranslation("mobile", trans, func(ut ut.Translator) error {
		return ut.Add("mobile", "{0}格式错误", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile", fe.Field(), fe.Field())
		return t
	}); err != nil {
		panic(err)
	}
}

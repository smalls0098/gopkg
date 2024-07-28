package validator

import (
	"reflect"
	"sync"

	"github.com/gin-gonic/gin/binding"
	enLocales "github.com/go-playground/locales/en"
	zhLocales "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

func Register() {
	binding.Validator = &CustomiValidator{}
}

var _ binding.StructValidator = &CustomiValidator{}

type CustomiValidator struct {
	once     sync.Once
	validate *validator.Validate
	Trans    ut.Translator
}

func (v *CustomiValidator) ValidateStruct(a interface{}) error {
	if kindOfData(a) != reflect.Struct {
		return nil
	}
	v.lazyInit()
	return v.validate.Struct(a)
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func (v *CustomiValidator) Engine() any {
	v.lazyInit()
	return v.validate
}

func (v *CustomiValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("validate")
		v.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			label := field.Tag.Get("label")
			if len(label) > 0 {
				return label
			}
			return field.Name
		})

		// 注册多语言
		zhT := zhLocales.New()
		enT := enLocales.New()
		uti := ut.New(enT, zhT, enT)
		v.Trans, _ = uti.GetTranslator("zh")
		if err := zhTranslations.RegisterDefaultTranslations(v.validate, v.Trans); err != nil {
			panic(err)
		}
		// 注册自定义验证方式
		registerMobile(v.validate, v.Trans)
	})
}

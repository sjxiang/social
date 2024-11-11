package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en" 
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// 验证器 
var validate *validator.Validate

// 翻译器 
var translator ut.Translator

func init() {

	// 这行代码创建了一个新的验证器实例，并启用了结构体字段的必填项检查。
	validate = validator.New(validator.WithRequiredStructEnabled())

	// 这行代码创建了一个新的翻译器实例，用于将验证错误信息翻译成英文。
	translator, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	// 这行代码注册了默认的英文翻译，以便在验证时能够使用这些翻译。
	en_translations.RegisterDefaultTranslations(validate, translator)

	// 这行代码注册了一个函数，用于获取结构体字段的 JSON 标签名称。如果 JSON 标签名称为 -，则返回空字符串。
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}


func Check(val any) error {
	// 这行代码使用 validate 验证器验证传入的 val 是否符合预定义的验证规则。
	if err := validate.Struct(val); err != nil {
		
		// 这行代码尝试将 err 变量转换为 validator.ValidationErrors 类型。
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, verror := range verrors {

			// 这行代码将每个验证错误信息转换为 FieldError 类型，并将其添加到 fields 切片中。
			fields = append(fields, FieldError{
				Field: verror.Field(),
				Err:   verror.Translate(translator),
			})
		}

		return fields
	}

	return nil
}

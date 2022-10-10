package validators

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"log"
	"reflect"
	"strings"
)

var (
	trans    ut.Translator
	jv       = &JsonValidator{}
	validate *validator.Validate
)

type Validation func() (string, validator.Func)
type Translation func() (string, validator.RegisterTranslationsFunc, validator.TranslationFunc)

type JsonValidator struct {
	Validation  []Validation
	Translation []Translation
}

func (jv *JsonValidator) LoadValidation(validations ...Validation) *JsonValidator {
	jv.Validation = validations
	return jv
}
func (jv *JsonValidator) LoadTranslation(translations ...Translation) *JsonValidator {
	jv.Translation = translations
	return jv
}

// Init 验证器初始化
func Init(locale string) {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
		// 注册翻译器
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		uni := ut.New(enT, zhT, enT)
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			log.Fatal("uni.GetTranslator(zh) failed")
		}

		// 注册翻译器
		var err error
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		if err != nil {
			log.Fatal("注册翻译器失败")
		}

		// 加载自定义验证器
		for _, v := range jv.Validation {
			tag, fn := v()
			if err := validate.RegisterValidation(tag, fn); err != nil {
				log.Fatal(fmt.Sprintf("注册%s验证器失败", tag))
			}
		}
		// 加载自定义翻译器
		for _, tran := range jv.Translation {
			tag, rfn, tfn := tran()
			if err := validate.RegisterTranslation(tag, trans, rfn, tfn); err != nil {
				return
			}
		}
	}
}

// RemoveTopStruct 移除多余的标签
func RemoveTopStruct(fields map[string]string) map[string]string {
	result := map[string]string{}
	for key, value := range fields {
		result[key[strings.Index(key, ".")+1:]] = value
	}
	return result
}

// InterceptError 拦截自定义Error
func InterceptError(er error) (ok bool, errors string) {
	if value, ok := er.(validator.ValidationErrors); ok {
		for i := 0; i < len(value); i++ {
			for _, value := range RemoveTopStruct(value.Translate(trans)) {
				errors = fmt.Sprintf("%s", value)
			}
		}
		return true, errors
	}
	return false, ""
}

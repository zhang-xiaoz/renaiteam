package tool

import (
	"fmt"
	"reflect"
	"strings"
	"web/models"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func InitVaildators() {
	//注册验证器
	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	_ = v.RegisterValidation("mobile", ValidatorMobile)
	// 	//下边代码是把自定义的验证器错误翻译由英语变中文//直接用就可以
	// 	v.RegisterTranslation("mobile", models.Trans, func(ut ut.Translator) error {
	// 		//fmt.Println("asasasassasasasasasasasasasas")
	// 		//ut.Add("mobile", "{0} 非法的手机号码!", true)
	// 		//fmt.Println("asasasassasasasasasasasasasas")
	// 		//fmt.Println(err)
	// 		return nil
	// 	}, func(utt ut.Translator, fe validator.FieldError) string {
	// 		t, _ := utt.T("mobile", fe.Field())
	// 		return t
	// 	})
	// }

	//初始化翻译
	InitTrans("zh")
}

// 翻译器
func InitTrans(locale string) (err error) {
	//修改gin框架中的validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//注册一个获取json的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" { //严谨为了json中-不处理
				return ""
			}
			return name
		})
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		//第一个参数是备用语言环境，后边的是应该支持的环境
		uni := ut.New(enT, zhT, enT)
		models.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("GetTranslator(%s)", locale)
		}
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(v, models.Trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, models.Trans)
		default:
			en_translations.RegisterDefaultTranslations(v, models.Trans)
		}
		return
	}
	return
}

// // 自定义验证器
// func ValidatorMobile(f1 validator.FieldLevel) bool {
// 	mobile := f1.Field().String() //转string
// 	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
// 	return ok
// }

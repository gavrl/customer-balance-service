package formatter

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
)

type JSONFormatter struct {
	translator ut.Translator
}

// NewJSONFormatter will create a new JSON formatter and register a custom tag
// name func to gin's validator
func NewJSONFormatter() (*JSONFormatter, error) {
	uni := ut.New(en.New())
	trans, _ := uni.GetTranslator("en")

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		//Verifier registration translator
		err := en_translations.RegisterDefaultTranslations(v, trans)
		if err != nil {
			return nil, err
		}

		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	return &JSONFormatter{trans}, nil
}

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func (frmt JSONFormatter) Descriptive(verr validator.ValidationErrors) []ValidationError {
	var errs []ValidationError

	for _, f := range verr {
		//err := f.ActualTag()
		err := f.Translate(frmt.translator)
		//if f.Param() != "" {
		//	err = fmt.Sprintf("%s=%s", err, f.Param())
		//}
		errs = append(errs, ValidationError{Field: f.Field(), Reason: err})
	}

	return errs
}

func (JSONFormatter) Simple(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs[f.Field()] = err
	}

	return errs
}

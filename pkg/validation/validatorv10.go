package validation

import (
	"errors"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
	"github.com/shandysiswandi/gostarter/pkg/strcase"
)

// V10Validator is a concrete implementation of the Validator interface.
type V10Validator struct {
	validate    *validator.Validate
	translators map[string]ut.Translator
}

type V10ValidationError map[string]string

func (vs V10ValidationError) Error() string {
	return "go-playground/validator/v10/custom-field-error"
}

func (vs V10ValidationError) Values() map[string]string {
	return vs
}

// NewV10Validator creates a new instance of V10Validator.
func NewV10Validator() (*V10Validator, error) {
	enLang := en.New()
	idLang := id.New()

	validate := validator.New(validator.WithRequiredStructEnabled())
	uni := ut.New(enLang, enLang, idLang)

	enTrans, _ := uni.GetTranslator("en")
	idTrans, _ := uni.GetTranslator("id")

	if err := enTranslations.RegisterDefaultTranslations(validate, enTrans); err != nil {
		return nil, err
	}
	if err := idTranslations.RegisterDefaultTranslations(validate, idTrans); err != nil {
		return nil, err
	}

	return &V10Validator{
		validate:    validate,
		translators: map[string]ut.Translator{"en": enTrans, "id": idTrans},
	}, nil
}

func AsV10Validator(err error) (V10ValidationError, bool) {
	var values V10ValidationError
	if !errors.As(err, &values) {
		return nil, false
	}

	return values, true
}

// Validate checks the given data against validation rules defined in its struct tags.
func (v *V10Validator) Validate(data any) error {
	if err := v.validate.Struct(data); err != nil {
		var validateErrs validator.ValidationErrors
		if !errors.As(err, &validateErrs) {
			return err
		}

		errV10 := make(V10ValidationError)
		for _, fe := range validateErrs {
			errV10[strcase.ToLowerSnake(fe.Field())] = fe.Translate(v.translators["en"])
		}

		return errV10
	}

	return nil
}

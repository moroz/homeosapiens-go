package helpers

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func LocalizeValidationErrors(l *i18n.Localizer, errs validation.Errors) validation.Errors {
	result := make(validation.Errors, len(errs))
	for field, err := range errs {
		if e, ok := errors.AsType[validation.Error](err); ok {
			translated, translateErr := l.Localize(&i18n.LocalizeConfig{
				MessageID:    "validation." + e.Code(),
				TemplateData: e.Params(),
				DefaultMessage: &i18n.Message{
					ID:    "validation." + e.Code(),
					Other: e.Error(),
				},
			})
			if translateErr != nil {
				result[field] = err
			} else {
				result[field] = errors.New(translated)
			}
		}
	}
	return result
}

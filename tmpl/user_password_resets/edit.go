package user_password_resets

import (
	"fmt"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Edit(ctx *types.CustomContext, token string, errors validation.Errors) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{
		ID: "user_password_resets.edit.title",
	})

	return layout.AuthLayout(ctx, title,
		Form(
			Class("mt-4 grid gap-4"),
			Method("POST"),
			Action(fmt.Sprintf("/reset-password/%s", token)),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "user_registrations.new.form.labels.password",
				}),
				Error:        errors,
				Name:         "password",
				ID:           "password",
				Type:         "password",
				Autocomplete: "new-password",
				Required:     true,
				HelperText: l.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "user_registrations.new.form.helper_text.password",
					TemplateData: map[string]string{
						"Min": strconv.Itoa(config.MinPasswordLength),
						"Max": strconv.Itoa(config.MaxPasswordLength),
					},
				}),
			}),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "user_registrations.new.form.labels.password_confirmation",
				}),
				Error:        errors,
				Name:         "password_confirmation",
				ID:           "password_confirmation",
				Type:         "password",
				Autocomplete: "new-password",
				Required:     true,
			}),

			Button(
				Type("submit"),
				Class("button"),
				Text(l.MustLocalizeMessage(&i18n.Message{
					ID: "user_password_resets.edit.submit",
				})),
			),
		),
	)
}

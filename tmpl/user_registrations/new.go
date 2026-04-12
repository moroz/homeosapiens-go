package userregistrations

import (
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

func New(ctx *types.CustomContext, params *types.RegisterUserParams, errors validation.Errors) Node {
	l := ctx.Localizer
	pageTitle := l.MustLocalizeMessage(&i18n.Message{
		ID: "user_registrations.new.title",
	})

	return layout.AuthLayout(ctx, pageTitle,
		components.GoogleButton(l.MustLocalizeMessage(&i18n.Message{
			ID: "user_registrations.new.sign_up_with_google",
		}), "", "w-full flex"),
		Hr(Class("my-4")),

		Form(
			Class("grid gap-4"),
			Method("POST"),
			Action("/sign-up"),
			Attr("novalidate", ""),
		
			Input(Type("hidden"), Name("locale"), Value(ctx.Language)),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "common.users.email",
				}),
				Value:        params.Email,
				Name:         "email",
				ID:           "email",
				Type:         "email",
				Autocomplete: "email",
				Required:     true,
				Error:        errors,
				Autofocus:    true,
			}),

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

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "common.users.given_name",
				}),
				Value:        params.GivenName,
				Error:        errors,
				Name:         "given_name",
				ID:           "given_name",
				Type:         "text",
				Autocomplete: "given-name",
				Required:     true,
			}),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "common.users.family_name",
				}),
				Value:        params.FamilyName,
				Error:        errors,
				Name:         "family_name",
				ID:           "family_name",
				Type:         "text",
				Autocomplete: "family-name",
				Required:     true,
			}),

			Button(
				Type("submit"),
				Class("button"),
				Text(l.MustLocalizeMessage(&i18n.Message{
					ID: "user_registrations.new.submit",
				})),
			),
		),

		Footer(
			Class("mt-4 text-center"),
			Raw(
				l.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "user_registrations.new.already_signed_up_html",
					TemplateData: map[string]string{
						"Path": "/sign-in",
					},
				}),
			),
		),
	)
}

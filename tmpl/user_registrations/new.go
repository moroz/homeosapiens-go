package userregistrations

import (
	"context"

	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"

	. "maragu.dev/gomponents/html"
)

func New(ctx context.Context) Node {
	l := ctx.Value("localizer").(*i18n.Localizer)
	pageTitle := l.MustLocalizeMessage(&i18n.Message{
		ID: "user_registrations.new.title",
	})

	return layout.AuthLayout(ctx, pageTitle,
		Form(
			Class("grid gap-4"),
			Method("POST"),
			Action("/sign-up"),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "user_registrations.new.form.labels.email",
				}),
				Name:         "email",
				ID:           "email",
				Type:         "email",
				Autocomplete: "email",
				Required:     true,
			}),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "user_registrations.new.form.labels.password",
				}),
				Name:         "password",
				ID:           "password",
				Type:         "password",
				Autocomplete: "new-password",
				Required:     true,
			}),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "user_registrations.new.form.labels.password_confirmation",
				}),
				Name:         "password_confirmation",
				ID:           "password_confirmation",
				Type:         "password",
				Autocomplete: "new-password",
				Required:     true,
			}),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "user_registrations.new.form.labels.given_name",
				}),
				Name:         "given_name",
				ID:           "given_name",
				Type:         "text",
				Autocomplete: "given-name",
				Required:     true,
			}),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "user_registrations.new.form.labels.family_name",
				}),
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

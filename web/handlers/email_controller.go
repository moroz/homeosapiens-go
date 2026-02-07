package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/tmpl/email"
)

type emailController struct{}

func EmailController() *emailController {
	return &emailController{}
}

func (c *emailController) Show(r *echo.Context) error {
	props := email.LayoutProps{
		Subject:     "Test subject",
		CompanyName: "Homeo sapiens",
		Heading:     "Confirm your email",
		Body:        "The quick brown fox jumps over the lazy dog.",
		ButtonText:  "Call to action",
		FooterText:  "&copy; 2024&ndash;2026 by Wydawnictwo Homeo Sapiens. All rights reserved.",
	}

	return email.LayoutTemplate.Execute(r.Response(), props)
}

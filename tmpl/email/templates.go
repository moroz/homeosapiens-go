package email

import (
	"embed"
	"html/template"
)

//go:embed *.html.tmpl
var templateFS embed.FS

type LayoutProps struct {
	Subject        string
	CompanyName    string
	Heading        string
	Body           template.HTML
	ButtonText     string
	FooterText     template.HTML
	ButtonURL      string
	UnsubscribeURL string
}

var LayoutTemplate = template.Must(template.ParseFS(templateFS, "layout.html.tmpl"))

package email

import (
	"embed"
	"html/template"

	"github.com/moroz/homeosapiens-go/types"
)

//go:embed *.html.tmpl
var templateFS embed.FS

type LayoutProps struct {
	LogoURL string
}

type OrderConfirmationEmailProps struct {
	*LayoutProps
	Order *types.OrderDTO
}

var LayoutTemplate = template.Must(template.ParseFS(templateFS, "layout.html.tmpl", "_header.html.tmpl", "_footer.html.tmpl"))

var OrderConfirmationTemplate = template.Must(template.ParseFS(templateFS, "layout.html.tmpl", "_header.html.tmpl", "_footer.html.tmpl", "order_confirmation.html.tmpl"))

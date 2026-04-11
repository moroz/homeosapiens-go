package handlers

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/email"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type emailController struct {
	db           queries.DBTX
	orderService *services.OrderRepository
}

func EmailController(db queries.DBTX) *emailController {
	return &emailController{
		orderService: services.NewOrderRepository(db),
		db:           db,
	}
}

func (cc *emailController) getLastOrder(c *echo.Context) (*types.OrderDTO, error) {
	lastID, err := queries.New(cc.db).GetLastOrderID(c.Request().Context())
	if err != nil {
		return nil, err
	}

	return cc.orderService.GetOrderDetails(c.Request().Context(), lastID)
}

func (cc *emailController) OrderConfirmation(c *echo.Context) error {
	order, err := cc.getLastOrder(c)
	if err != nil {
		return err
	}

	ctx := helpers.GetRequestContext(c)

	subject := ctx.Localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    "emails.order_confirmation.subject",
		TemplateData: order,
	})

	props := email.OrderEmailProps{
		LayoutProps: &email.LayoutProps{
			Title:     subject,
			Localizer: ctx.Localizer,
			Language:  ctx.Language,
		},
		Order: order,
	}

	if err := email.OrderConfirmationTemplate.Execute(c.Response(), props); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (cc *emailController) PaymentConfirmation(c *echo.Context) error {
	order, err := cc.getLastOrder(c)
	if err != nil {
		return err
	}

	ctx := helpers.GetRequestContext(c)

	subject := ctx.Localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    "emails.payment_confirmation.subject",
		TemplateData: order,
	})

	props := email.OrderEmailProps{
		LayoutProps: &email.LayoutProps{
			Title:     subject,
			Localizer: ctx.Localizer,
			Language:  ctx.Language,
		},
		Order: order,
	}

	return email.PaymentConfirmationTemplate.Execute(c.Response(), props)
}

func (cc *emailController) UserEmailVerification(c *echo.Context) error {
	user, err := queries.New(cc.db).GetUserByEmail(c.Request().Context(), crypto.HashEmail("karol@moroz.dev"))
	if err != nil {
		return err
	}

	token, err := services.NewUserTokenService(cc.db).IssueEmailVerificationTokenForUser(c.Request().Context(), user, config.EmailVerificationTokenValidity)

	ctx := helpers.GetRequestContext(c)

	subject := ctx.Localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    "emails.user_email_verification.subject",
		TemplateData: user,
	})

	props := email.UserEmailVerificationEmailProps{
		LayoutProps: &email.LayoutProps{
			Title:     subject,
			Language:  ctx.Language,
			Localizer: ctx.Localizer,
		},
		UserToken: token,
	}

	return email.UserEmailVerificationTemplate.Execute(c.Response(), props)
}

package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/email_verifications"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type emailVerificationController struct {
	db  queries.DBTX
	srv *services.EmailVerificationService
}

func EmailVerificationController(db queries.DBTX) *emailVerificationController {
	return &emailVerificationController{
		db:  db,
		srv: services.NewEmailVerificationService(db),
	}
}

func (cc *emailVerificationController) New(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	email := c.QueryParam("email")
	return email_verifications.New(ctx, email, "").Render(c.Response())
}

func (cc *emailVerificationController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	var params types.ResendEmailVerificationTokenParams
	if err := c.Bind(&params); err != nil {
		return echo.ErrBadRequest
	}

	redirectTo := c.Request().Referer()
	if redirectTo == "" {
		redirectTo = "/sign-in"
	}

	l := ctx.Localizer

	_, err := cc.srv.MaybeResendVerificationEmail(c.Request().Context(), params.Email)
	if errors.Is(err, services.ErrRateLimited) {
		ctx.PutFlash("error", l.MustLocalizeMessage(&i18n.Message{
			ID: "email_verifications.create.rate_limited",
		}))
	} else if err != nil {
		log.Print(err)
		isGmail := strings.HasSuffix(strings.TrimSpace(strings.ToLower(params.Email)), "@gmail.com")
		msgKey := "email_verifications.create.error_not_gmail"
		if isGmail {
			msgKey = "email_verifications.create.error_is_gmail"
		}

		ctx.PutFlash("error", l.MustLocalizeMessage(&i18n.Message{
			ID: msgKey,
		}))
	} else {
		ctx.PutHTMLFlash("success", l.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "email_verifications.create.success_flash",
			TemplateData: map[string]string{
				"Email": params.Email,
			},
		}))
	}

	if err := ctx.SaveSession(c.Response()); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, redirectTo)
}

func (cc *emailVerificationController) Verify(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	token := c.QueryParam("token")
	if token == "" {
		return echo.ErrBadRequest
	}

	user, err := cc.srv.VerifyEmailAddress(c.Request().Context(), token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	ctx.PutFlash("success", ctx.Localizer.MustLocalizeMessage(&i18n.Message{
		ID: "email_verifications.verify.success",
	}))

	return signUserIn(c, cc.db, user)
}

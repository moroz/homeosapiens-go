package handlers

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/services"
	userregistrations "github.com/moroz/homeosapiens-go/tmpl/user_registrations"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type userRegistrationController struct {
	db  queries.DBTX
	srv *services.UserRegistrationService
}

func UserRegistrationController(db queries.DBTX) *userRegistrationController {
	return &userRegistrationController{
		db:  db,
		srv: services.NewUserRegistrationService(db),
	}
}

func (cc *userRegistrationController) New(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	return userregistrations.New(ctx, &types.RegisterUserParams{}, make(validation.Errors)).Render(c.Response())
}

func (cc *userRegistrationController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	var params types.RegisterUserParams
	if err := c.Bind(&params); err != nil {
		log.Print(err)
		return echo.ErrBadRequest
	}

	user, err := cc.srv.RegisterUser(c.Request().Context(), &params)
	if err, ok := errors.AsType[validation.Errors](err); ok {
		validationErrors := helpers.LocalizeValidationErrors(ctx.Localizer, err)
		return userregistrations.New(ctx, &params, validationErrors).Render(c.Response())
	}
	if err != nil {
		ctx.PutFlash("error", err.Error())
		if err := ctx.SaveSession(c.Response()); err != nil {
			log.Print(err)
		}
		return userregistrations.New(ctx, &params, nil).Render(c.Response())
	}

	redirectTo := config.PublicUrl + "/user-registrations/success?token=" + user.UserTokenDTO.EncodeToken()

	return c.Redirect(http.StatusFound, redirectTo)
}

func (cc *userRegistrationController) Success(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	param := c.QueryParam("token")
	token, err := base64.RawURLEncoding.DecodeString(param)
	if param == "" || err != nil {
		return echo.ErrBadRequest
	}

	user, err := queries.New(cc.db).FindUserByRegistrationToken(c.Request().Context(), crypto.HashUserToken(token))
	if err != nil {
		log.Print(err)
	}

	return userregistrations.Success(ctx, user, param).Render(c.Response())
}

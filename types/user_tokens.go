package types

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type UserTokenDTO struct {
	*queries.UserToken
	*queries.User
	PlaintextToken []byte
}

func (u *UserTokenDTO) EmailRecipient() string {
	return fmt.Sprintf("%s %s <%s>", u.GivenName.String(), u.FamilyName.String(), u.Email.String())
}

func (u *UserTokenDTO) EncodeToken() string {
	return base64.RawURLEncoding.EncodeToString(u.PlaintextToken)
}

func (u *UserTokenDTO) VerifyEmailPath() string {
	qs := url.Values{"token": {u.EncodeToken()}}.Encode()
	return "/verify-email?" + qs
}

func (u *UserTokenDTO) VerifyEmailURL() string {
	return config.PublicUrl + u.VerifyEmailPath()
}

package helpers

import (
	"testing"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/sessions"
	"github.com/stretchr/testify/assert"
)

func TestGetRedirectUrl(t *testing.T) {
	cases := []struct {
		stored string
		want   string
	}{
		{"", "/"},
		{"/profile", "/profile"},
		{"/videos/group/slug?tab=1", "/videos/group/slug?tab=1"},
		{"//host/favicon.ico", "/"}, // protocol-relative (the favicon bug)
		{"//evil.com", "/"},         // open redirect
		{"/\\evil.com", "/"},        // backslash protocol-relative
		{"https://evil.com", "/"},   // absolute URL
		{"http://host/x", "/"},      // absolute URL
		{"favicon.ico", "/"},        // not an absolute path
	}

	for _, tc := range cases {
		ctx := &types.CustomContext{Session: sessions.Payload{}}
		if tc.stored != "" {
			ctx.Session[config.RedirectBackUrlSessionKey] = tc.stored
		}

		got := GetRedirectUrl(ctx)

		assert.Equalf(t, tc.want, got, "GetRedirectUrl(%q)", tc.stored)
		assert.NotContainsf(t, ctx.Session, config.RedirectBackUrlSessionKey,
			"GetRedirectUrl(%q) should consume the stored value", tc.stored)
	}
}

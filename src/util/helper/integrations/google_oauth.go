package integrations

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	UserInfoURL          = "https://www.googleapis.com/oauth2/v2/userinfo?access="
	Provider             = "https://accounts.google.com"
	ScopesURLUserInfo    = "https://www.googleapis.com/auth/userinfo.email"
	ScopesURLUserProfile = "https://www.googleapis.com/auth/userinfo.profile"
)

var (
	SSOSignUp          *oauth2.Config
	ClientIDSignUp     = os.Getenv("GOOGLE_CLIENT_ID")
	ClientSecretSignUp = os.Getenv("GOOGLE_CLIENT_SECRET")
	RedirectURLSignUp  = os.Getenv("GOOGLE_REDIRECT_URL")
	RandomString       = "123qwerty"
)

func init() {
	SSOSignUp = initAuthConfig(ClientIDSignUp, ClientSecretSignUp, RedirectURLSignUp)
}

func initAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{ScopesURLUserInfo, ScopesURLUserProfile},
		Endpoint:     google.Endpoint,
	}
}

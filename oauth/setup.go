package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
	"tpbt/services"
)

var (
	Configuration *oauth2.Config
)

func Initialize(prv *services.Provider) {
	Configuration = &oauth2.Config {
		RedirectURL: "http://localhost:8081/#/auth/callback",
		ClientID: prv.Config.ClientID,
		ClientSecret: prv.Config.ClientSecret,
		Endpoint: twitch.Endpoint,
	}
}

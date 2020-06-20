package helix

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"tpbt/bterrors"
	"tpbt/models"
	"tpbt/services"
)

var (
	HttpClient = &http.Client{}
)

type TwitchUser struct {
	ID          string `json:"id"`
	Name        string `json:"login"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"profile_image_url"`
}

func FetchUser(provider *services.Provider, player *models.BTStreamer, continuing bool) error {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Client-ID", provider.Config.ClientID)
	req.Header.Set("Authorization", "Bearer " + player.AccessToken)

	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		txt, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		response := struct {
			Data []TwitchUser `json:"data"`
		}{}

		err = json.Unmarshal(txt, &response)
		if err != nil {
			return err
		}

		if len(response.Data) == 0 {
			return bterrors.ErrorNoUser
		}

		player.TwitchID = response.Data[0].ID
		player.TwitchName = response.Data[0].Name
		player.DisplayName = response.Data[0].DisplayName
		player.Picture = response.Data[0].Picture

		return nil
	} else if resp.StatusCode == 401 && continuing {
		// The token was refused! Let's try to refresh
		if !RefreshToken(provider, player) {
			return bterrors.ErrorCantRefresh
		}

		return FetchUser(provider, player, false)
	}

	return bterrors.ErrorAPI
}

func RefreshToken(provider *services.Provider, player *models.BTStreamer) bool {

	body := url.Values{}
	body.Set("grant_type", "refresh_token")
	body.Set("client_id", provider.Config.ClientID)
	body.Set("client_secret", provider.Config.ClientSecret)
	body.Set("refresh_token", player.RefreshToken)

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token?grand", strings.NewReader(body.Encode()))
	if err != nil {
		return false
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false
	}

	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var refreshed struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Scope string `json:"scope"`
	}

	err = json.Unmarshal(bodyResp, &refreshed)
	if err != nil {
		return false
	}

	player.AccessToken = refreshed.AccessToken
	player.RefreshToken = refreshed.RefreshToken

	// Save the new token in DB
	_, err = provider.DB.Exec("UPDATE twitch_users SET access_token = $1, refresh_token = $2 WHERE tw_id = $3", refreshed.AccessToken, refreshed.RefreshToken, player.TwitchID)
	if err != nil {
		return false
	}

	return true
}

package dal

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"tpbt/bterrors"
	"tpbt/models"
	"tpbt/services"
)

func FetchChannelsToJoin(prv *services.Provider) ([]string, error) {
	rows, err := prv.DB.Queryx("SELECT tw_name FROM TWITCH_USERS WHERE bot_in_channel = true")
	if err != nil {
		return []string{}, err
	}

	var users []string

	for rows.Next() {
		var u string

		err := rows.Scan(&u)
		if err != nil {
			continue
		}

		users = append(users, u)
	}

	return users, nil
}

func FetchUserFromToken(prv *services.Provider, token string) (*models.BTStreamer, error) {
	row := prv.DB.QueryRowx("SELECT tw_id, app_token, access_token, refresh_token, bot_in_channel, grp.ROLE as ROLE FROM TWITCH_USERS tw INNER JOIN GROUPS grp ON tw.GRP_ID = grp.GRP_ID WHERE APP_TOKEN = $1", token)

	u := &models.BTStreamer{
		MutexWS: &sync.Mutex{},
	}
	err := row.StructScan(u)
	if err != nil {
		return nil, bterrors.ErrorTokenNotFound
	}

	return u, nil
}

// FindOrCreateUser gets a user which at least contains a TwitchID and either pull the rest of the data from the DB or insert it in it
func FindOrCreateUser(prv *services.Provider, user *models.BTStreamer) error {
	uuidV4, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("Can't generate uuid!")
		return err
	}

	row := prv.DB.QueryRowx(`
		SELECT tw_id, tw_name, app_token, access_token, refresh_token, bot_in_channel, grp.ROLE as ROLE 
		FROM TWITCH_USERS tw
			 INNER JOIN GROUPS grp ON tw.GRP_ID = grp.GRP_ID
		WHERE TW_ID = $1
	`, user.TwitchID)
	err = row.StructScan(user)

	user.AppToken = uuidV4.String()

	if err != nil {
		if err == sql.ErrNoRows {
			// The user does not exists so we create it

			_, err = prv.DB.Exec("INSERT INTO twitch_users (tw_id, tw_name, app_token, access_token, refresh_token, bot_in_channel, grp_id) VALUES ($1, $2, $3, $4, $5, 'f', 1)", user.TwitchID, user.TwitchName, user.AppToken, user.AccessToken, user.RefreshToken)
			if err != nil {
				return err
			}

		}
		return err
	} else {
		// Update the stored token
		_, err = prv.DB.Exec("UPDATE twitch_users SET app_token = $1, tw_name = $2, access_token = $3, refresh_token = $4 WHERE tw_id = $5", user.AppToken, user.TwitchName, user.AccessToken, user.RefreshToken, user.TwitchID)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpsertUser inserts the user if it does not exists, update his username if it does
func UpsertUser(prv *services.Provider, user *models.BTPlayer) {
	rq := `
		INSERT INTO VIEWER (tw_id, last_known_username)
		VALUES ($1, $2)
		ON CONFLICT (tw_id) DO 
		UPDATE SET last_known_username = $2
	`

	_, err := prv.DB.Exec(rq, user.ID, user.Name)
	if err != nil {
		fmt.Println(err)
	}
}

// FetchPlayerForChannel returns a viewer for the given channel if he exists
func FetchPlayerForChannel(prv *services.Provider, streamer string, player string) (*models.BTPlayer, error){
	row := prv.DB.QueryRowx(`
			SELECT p.tw_id, p.last_known_username, COALESCE(sb.score, 0) score
			FROM VIEWER p
				 LEFT JOIN SCOREBOARD sb ON sb.viewer_id = p.tw_id
			WHERE p.tw_id = $1
				  AND (sb.streamer_id is null OR sb.streamer_id = $2)
	`, player, streamer)

	pl := &models.BTPlayer{}
	err := row.StructScan(pl)

	return pl, err
}

func UpdateBotInChannel(prv *services.Provider, user *models.BTStreamer) error {
	_, err := prv.DB.Exec("UPDATE TWITCH_USERS SET BOT_IN_CHANNEL = $1 WHERE TW_ID = $2", user.IsBotInChannel, user.TwitchID)
	return err
}
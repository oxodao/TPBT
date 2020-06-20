package dal

import (
	"database/sql"
	"fmt"
	"tpbt/models"
	"tpbt/services"
)

type Score struct {
	Name  string `db:"name"`
	Score int    `db:"score"`
}

func FetchScore (prv *services.Provider, streamer, user string) (*Score, error) {
	row := prv.DB.QueryRowx(`
		SELECT  str.tw_name as name,
				sb.score AS score
		FROM scoreboard sb
			 INNER JOIN twitch_users str ON str.tw_id = sb.STREAMER_ID
		WHERE sb.STREAMER_ID = $1
          AND sb.VIEWER_ID = $2`, streamer, user)

	score := Score{}
	err := row.StructScan(&score)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &score, err
}

func UpsertScore(prv *services.Provider, streamer *models.BTStreamer, player *models.BTPlayer) (error) {
	var score int

	// @TODO: Score customizable for each game
	if player.TitleFound {
		score = score + 50
	}

	if player.ArtistFound {
		score = score + 50
	}

	UpsertUser(prv, player)

	rq := `
		INSERT INTO scoreboard (viewer_id, streamer_id, score)
		VALUES ($1, $2, $3)
		ON CONFLICT (viewer_id, streamer_id) DO
		UPDATE SET score = $3
	`

	_, err := prv.DB.Exec(rq, player.ID, streamer.TwitchID, score)

	return err
}

func FetchScoreboard(prv *services.Provider, streamer *models.BTStreamer) ([]*Score, error) {
	rq := ` SELECT  viewer.last_known_username as name,
					sb.score AS score
			FROM scoreboard sb
				 INNER JOIN viewer ON viewer.tw_id = sb.VIEWER_ID
			WHERE sb.STREAMER_ID = $1
			ORDER BY sb.score DESC
			LIMIT 100`

	rows, err := prv.DB.Queryx(rq, streamer.TwitchID)
	if err != nil {
		return nil, err
	}

	leaderboard := []*Score{}
	for rows.Next() {
		sc := &Score{}
		err = rows.StructScan(sc)
		if err != nil {
			fmt.Println(err)
			continue
		}
		leaderboard = append(leaderboard, sc)
	}

	return leaderboard, nil
}

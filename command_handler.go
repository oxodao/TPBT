package main

import (
	"encoding/json"
	"fmt"
	"tpbt/dal"
	"tpbt/models"
	"tpbt/services"
)

func ProcessToggleBot(prv *services.Provider, user *models.BTStreamer, cmd *models.Command) {
	user.IsBotInChannel = !user.IsBotInChannel

	err := dal.UpdateBotInChannel(prv, user)
	if err != nil {
		fmt.Println(err)
	}

	err = SendCommand(user, "SET_USER", user)
	if err != nil {
		fmt.Println(err)
	}

	if user.IsBotInChannel {
		fmt.Printf("[%v] Joining channel\n", user.TwitchName)
		prv.Twitch.Join(user.TwitchName)
		prv.Twitch.Say(user.TwitchName, "Bonsoir bonsoir, je suis de retour!")
	} else {
		fmt.Printf("[%v] Leaving channel\n", user.TwitchName)
		prv.Twitch.Say(user.TwitchName, "J'y vais moi, la bise!")
		prv.Twitch.Depart(user.TwitchName)
	}
}

func ProcessToggleTurn(prv *services.Provider, user *models.BTStreamer, cmd *models.Command) {
	if user.Game.Running {
		user.Game.Running = false

		scores, err := dal.FetchScoreboard(prv, user)
		if err != nil {
			fmt.Println("Something went wrong fetching scores")
		} else {
			_ = SendCommand(user, "SET_LEADERBOARD", scores)
		}

		lbTurn := "La manche est terminé. La musique était " + user.Game.Title + " de " + user.Game.Artist + "."
		length := len(user.Game.Players)
		if  length > 0 {
			if length == 1 {
				lbTurn = lbTurn + " Le gagnant est " + user.Game.Players[0].Name
			} else if length == 2 {
				lbTurn = lbTurn + " Voici le podium: #1 " + user.Game.Players[0].Name + ", #2 " + user.Game.Players[1].Name
			} else {
				lbTurn = lbTurn + " Voici le podium: #1 " + user.Game.Players[0].Name + ", #2 " + user.Game.Players[1].Name + ", #3 " + user.Game.Players[2].Name
			}
		}

		prv.Twitch.Say(user.TwitchName, lbTurn)
	} else {
		game := &models.Game{}
		err := json.Unmarshal([]byte(cmd.Arguments), game)
		if err != nil {
			fmt.Println(err)
		}

		game.Running = true
		game.Players = []*models.BTPlayer{}
		user.Game = game

		_ = SendCommand(user, "SET_FOUND", game.Players)

		prv.Twitch.Say(user.TwitchName, "C'est parti ! Écoutez la musique et devinez le titre et l'artiste.")
	}

	_ = SendCommand(user, "SET_GAME", user.Game)
}

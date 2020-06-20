package main

import (
	"database/sql"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"regexp"
	"strings"
	"tpbt/dal"
	"tpbt/models"
	"tpbt/services"
)

var (
	RegExAnswers *regexp.Regexp
)

func ConfigureBot(prv *services.Provider) {
	prv.Twitch.OnPrivateMessage(OnPrivateMessage(prv))

	var err error
	RegExAnswers, err = regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic("Ton pc est autiste")
	}

	channelsToJoin, err := dal.FetchChannelsToJoin(prv)
	if err != nil {
		panic(err)
	}

	for _, c := range channelsToJoin {
		prv.Twitch.Join(c)
	}

	err = prv.Twitch.Connect()
	if err != nil {
		panic(err)
	}
}

func OnPrivateMessage(prv *services.Provider) func(message twitch.PrivateMessage) {
	return func (message twitch.PrivateMessage) {
		msg := message.Message
		fmt.Println(msg)
		if strings.HasPrefix(msg, "!tpbt") {
			if len(msg) == 5 { // If there are nothing else than the command name
				// Send his score to the player
				score, err := dal.FetchScore(prv, message.RoomID, message.User.ID)
				if err != nil {
					fmt.Printf("[%v] Something went wrong loading score for user %v.\n", message.ID, message.User.DisplayName)
				}

				fmt.Println("Sending score to user")
				prv.Twitch.Whisper(message.User.Name, fmt.Sprintf("[%v] Score: %v.", score.Name, score.Score))
			}
		} else if streamer, ok := Users[message.RoomID]; ok {
			if streamer.Game.Running {
				// If the streamer is connected and the game is running (i.e. the turn is started and the timer counts down)

				var user *models.BTPlayer
				ok := false
				for _, p := range streamer.Game.Players {
					if p.ID == message.User.ID {
						user = p
						ok = true
					}
				}

				// We check if the user is already presents in the current answers
				if !ok {
					// If this is not the case, we try to get it from the database
					player, err := dal.FetchPlayerForChannel(prv, message.RoomID, message.User.ID)
					if err != nil {
						if err != sql.ErrNoRows {
							fmt.Println("Something went wrong searching for the user!")
							fmt.Println(err)

							return
						}
						player.ID = message.User.ID
						player.Name = message.User.DisplayName
					}
					streamer.Game.Players = append(streamer.Game.Players, player)
					user = player
				}

				user.Game = streamer.Game

				// If he has found everything, he's probably simply talking into the chat.
				// No need to process his answer
				if !user.TitleFound || !user.ArtistFound {
					CheckAnswer(prv, streamer, user, message.Message)
				}
			}
		}
	}
}

func normalizedContains(s1, s2 string) bool {
	tmpS1 := strings.ToLower(s1)
	tmpS1 = RegExAnswers.ReplaceAllString(tmpS1, "")

	tmpS2 := strings.ToLower(s2)
	tmpS2 = RegExAnswers.ReplaceAllString(tmpS2, "")

	return strings.Contains(tmpS1, tmpS2)
}

type PlayerFound struct {
	Name        string
	FoundArtist bool
	FoundTitle  bool
}

func CheckAnswer(prv *services.Provider, streamer *models.BTStreamer, player *models.BTPlayer, msg string) {
	foundSomething := false

	if normalizedContains(msg, streamer.Game.Title) {
		// Just to prevent re-saving into DB if the user answers multiple times
		if !player.TitleFound {
			foundSomething = true
		}
		player.SetFoundTitle()
	}

	if normalizedContains(msg, streamer.Game.Artist) {
		// Just to prevent re-saving into DB if the user answers multiple times
		if !player.ArtistFound {
			foundSomething = true
		}
		player.SetFoundArtist()
	}

	if foundSomething {
		if player.TitleFound && !player.ArtistFound{
			prv.Twitch.Whisper(player.Name, fmt.Sprintf("[%v] Bravo, le titre est bien %v", streamer.DisplayName, streamer.Game.Title))
			fmt.Printf("[%v] %v found title.\n", streamer.DisplayName, player.Name)
		}

		if !player.TitleFound && player.ArtistFound {
			prv.Twitch.Whisper(player.Name, fmt.Sprintf("[%v] Bravo, l'artiste est bien %v", streamer.DisplayName, streamer.Game.Artist))
			fmt.Printf("[%v] %v found artist.\n", streamer.DisplayName, player.Name)
		}

		if player.TitleFound && player.ArtistFound {
			prv.Twitch.Whisper(player.Name, fmt.Sprintf("[%v] Bien joué, il s'agissait effectivement de %v", streamer.DisplayName, streamer.Game.Title + " - " + streamer.Game.Artist))
			fmt.Printf("[%v] %v found both.\n", streamer.DisplayName, player.Name)
		}

		// Save in DB
		err := dal.UpsertScore(prv, streamer, player)
		if err != nil {
			fmt.Println("Score not saved!")
			fmt.Println(err)
		}

		_ = SendCommand(streamer, "SET_FOUND", streamer.Game.Players)
	}
}

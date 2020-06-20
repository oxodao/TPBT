package main

import (
	"encoding/json"
	"fmt"
	"tpbt/models"
	"tpbt/services"
)

var (
	Users = make(map[string]*models.BTStreamer)
)

func ListenWS(prv *services.Provider, user *models.BTStreamer) {
	user.IsConnected = true
	Users[user.TwitchID] = user

	user.Websocket.SetCloseHandler(func(code int, text string) error {
		if _, ok := Users[user.TwitchID]; ok {
			delete(Users, user.TwitchID)
		}

		user.IsConnected = false
		user.Websocket.Close()
		return nil
	})

	user.Game = &models.Game{}

	go ReceiveCommands(prv, user)
}

func ReceiveCommands(prv *services.Provider, user *models.BTStreamer) {
	for user.IsConnected {
		_, msg, err := user.Websocket.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}

		cmd := &models.Command{}
		err = json.Unmarshal(msg, cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ExecuteCommand(prv, user, cmd)
	}
}

func ExecuteCommand(prv *services.Provider, user *models.BTStreamer, cmd *models.Command) {
	switch cmd.Command {
	case "PING":
		break
	case "TOGGLE_BOT":
		ProcessToggleBot(prv, user, cmd)
		break
	case "TOGGLE_TURN":
		ProcessToggleTurn(prv, user, cmd)
		break
	default:
		fmt.Println("Unhandled command: " + cmd.Command)
		break
	}
}

func SendCommand(u *models.BTStreamer, cmd string, payload interface{}) error {
	content, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	u.MutexWS.Lock()
	err = u.Websocket.WriteJSON(models.Command{
		Command:   cmd,
		Arguments: string(content),
	})
	u.MutexWS.Unlock()

	return err
}

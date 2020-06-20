package models

import (
	"github.com/gorilla/websocket"
	"sync"
)

type BTStreamer struct {
	TwitchID     string `json:"-" db:"tw_id"`
	TwitchName   string `json:"-" db:"tw_name"`
	AccessToken  string `json:"-" db:"access_token"`
	RefreshToken string `json:"-" db:"refresh_token"`

	IsBotInChannel bool `json:"IsBotInChannel" db:"bot_in_channel"`

	AppToken     string `json:"AppToken" db:"app_token"`
	DisplayName  string `json:"DisplayName" db:"-"`
	Picture      string `json:"Picture" db:"-"`
	Role         string `json:"Role" db:"role"`

	Game *Game `json:"Game" db:"-"`
	Leaderboard []*BTPlayer `json:"Leaderboard" db:"-"`

	IsConnected  bool            `json:"-" db:"-"`
	Websocket    *websocket.Conn `json:"-" db:"-"`
	MutexWS      *sync.Mutex     `json:"-" db:"-"`
}


package services

import (
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/jmoiron/sqlx"
	"tpbt/config"
)

type Provider struct {
	Config *config.Configuration
	DB *sqlx.DB
	Twitch *twitch.Client
}

func NewProvider(cfg *config.Configuration) (*Provider, error) {
	prv := &Provider {
		Config: cfg,
	}

	dbc := cfg.Database
	sqlConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbc.Host, dbc.Port, dbc.Username, dbc.Password, dbc.Database)

	/**
	 * Database
	 */
	db, err := sqlx.Connect("postgres", sqlConnection)
	if err != nil {
		return nil, err
	}
	prv.DB = db

	/**
	 * IRC Chat
	 */
	prv.Twitch = twitch.NewClient(cfg.Bot.Username, cfg.Bot.Secret)

	return prv, nil
}
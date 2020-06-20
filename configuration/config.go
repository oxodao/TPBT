package configuration

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Port string
	ClientID string
	ClientSecret string
	Database struct {
		Host string
		Port string
		Username string
		Password string
		Database string
	}
	Bot struct {
		Username string
		Secret string
	}
}

func LoadConfiguration() (*Configuration, error) {
	cfg := Configuration{}

	txt, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return &cfg, err
	}

	err = json.Unmarshal(txt, &cfg)
	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}

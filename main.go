package main

import (
	"fmt"
	"time"
	"tpbt/configuration"
	"tpbt/oauth"
	"tpbt/services"

	_ "github.com/lib/pq"
)

const (
	tpbtVersion = "1.0"
	tpbtAuthor = "Oxodao"
)

func main() {
	fmt.Printf("TwitchPlaysBlindTest [v%v] by %v\n", tpbtVersion, tpbtAuthor)
	cfg, err := configuration.LoadConfiguration()
	if err != nil {
		panic(err)
	}

	prv, err := services.NewProvider(cfg)
	if err != nil {
		panic(err)
	}

	oauth.Initialize(prv)

	go ConfigureBot(prv)
	go RunServer(prv)

	for {
		time.Sleep(1 * time.Hour)
	}
}

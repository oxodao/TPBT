package main

import (
	"time"
	"tpbt/models"
)

func RunCountdown(g *models.Game, u *models.BTStreamer, callback func()) {
	doTheCount := true

	for doTheCount {
		doTheCount = g.Running && g.TimeLeft > 1

		g.TimeLeft -= 1
		_ = SendCommand(u, "COUNTDOWN_TICK", nil)

		time.Sleep(1 * time.Second)
	}

	callback()
}

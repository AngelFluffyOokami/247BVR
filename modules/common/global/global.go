package global

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tjgq/broadcast"
)

// Session todo
var Session *discordgo.Session

// Config todo
var Config Data

var PauseThreads = make(chan bool)
var paused = false
var PauseAll = broadcast.New(1)

func PauseExec() {
	for {
		p := <-PauseThreads
		if p {
			if !paused {
				paused = true
				PauseAll.Send(true)
			}
		} else if !p {
			if paused {
				paused = false
				PauseAll.Send(false)
			}
		}

	}
}

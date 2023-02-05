package common

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tjarratt/babble"
	"gorm.io/gorm"
)

func BabbleWords() string {
	var wordlist []string
	wordlist = append(wordlist, "grudge", "linear", "burial", "latest", "screen", "desert", "expose", "endure", "estate", "master", "refund", "throat", "effort", "pepper", "budget", "revive", "breast", "school", "flower", "ladder", "chorus", "wonder", "cheese", "sticky", "spread", "tumble", "vacuum", "flavor", "suntan", "mutter", "center", "punish", "resort", "hunter", "galaxy", "charge", "depend", "cotton", "shiver", "afford", "agenda", "timber", "morale", "behave", "camera", "expand", "carbon", "dollar", "latest", "mature", "mobile", "injury", "ensure", "barrel", "finish", "rhythm", "crutch", "museum", "lesson", "follow", "please", "safety", "modest", "remind", "reader", "demand", "ethics", "pledge", "accept", "ballot", "doctor", "gutter", "planet", "launch", "makeup", "freeze", "acquit", "colony", "rescue", "defend", "facade", "vision", "honest", "retire", "arrest", "banner", "thesis", "weight", "turkey", "worker", "column", "ignite", "facade", "ribbon", "bloody", "sacred", "inside", "dilute", "gallon", "theory", "behead", "proper", "chance", "single", "object", "temple", "modest", "likely", "adjust", "pastel", "attack", "market", "bishop", "belong", "effort", "rotate", "senior", "infect", "locate", "secure", "earwax", "normal", "flower", "prayer", "endure", "injury", "avenue", "family", "desert", "packet", "series", "tiptoe", "tumble", "harass", "spider", "output", "mutter", "church", "glance", "throne", "salmon", "option", "apathy", "cancer", "labour", "stroke", "dinner", "lounge", "gallon", "mobile", "bubble", "trance", "matrix", "ground", "escape", "defeat", "effect", "acquit", "square", "bitter", "excuse", "review", "normal", "formal", "player", "quaint", "belief", "critic", "accent", "empire", "junior", "lesson", "tongue", "voyage", "basket", "launch", "mosaic", "column", "margin", "source", "spirit", "cherry", "height", "bother", "deadly", "marble", "virtue", "devote", "mosque", "morale", "likely", "branch", "offend", "family", "script", "medium", "course", "theory", "weight", "winner")
	babbler := babble.NewBabbler()
	babbler.Count = 6
	babbler.Words = wordlist
	key := babbler.Babble()
	return key
}

const (
	LogError    = "ERR"
	LogWarning  = "WARN"
	LogInfo     = "INFO"
	LogUpdate   = "UPDATE"
	LogFeedback = "FEEDBACK"
)

func GetGuildName(GID string) string {
	s := Session
	guild, err := s.Guild(GID)
	if err != nil {
		return "Undefined. " + GID
	} else {
		return guild.Name + " " + GID
	}
}

func GetGuildOwnerName(GID string) string {
	s := Session
	g, err := s.Guild(GID)
	if err != nil {
		return "Undefined."
	}

	user, err := s.User(g.OwnerID)
	if err != nil {
		return "Undefined. " + g.OwnerID
	} else {
		return user.Username + "#" + user.Discriminator + " " + g.OwnerID
	}
}

func LogEvent(message string, level string) {

	config := Config
	s := Session
	// create the log entry
	entry := LogEntry{
		Time:    time.Now(),
		Message: message,
		Level:   level,
	}

	// open the log file
	logFile, err := os.OpenFile("opossum.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// marshal the log entry
	entryJSON, err := json.MarshalIndent(entry, "\n", "\n")
	if err != nil {
		panic(err)
	}

	// write the log entry to the file
	if _, err := logFile.Write(entryJSON); err != nil {
		panic(err)
	}
	if err := logFile.Sync(); err != nil {
		panic(err)
	}
	switch level {
	case "INFO":
		s.ChannelMessageSend(config.InfoChannel, entry.Level+": \n"+entry.Message+"\n"+fmt.Sprint(entry.Time))
	case "WARN":
		s.ChannelMessageSend(config.WarnChannel, entry.Level+": \n"+entry.Message+"\n"+fmt.Sprint(entry.Time))
	case "ERR":
		s.ChannelMessageSend(config.ErrChannel, entry.Level+": \n"+entry.Message+"\n"+fmt.Sprint(entry.Time))
	case "UPDATE":
		s.ChannelMessageSend(config.UpdateChannel, entry.Level+": \n"+entry.Message+"\n"+fmt.Sprint(entry.Time))
	}
}

var Config Data

var DefaultID = "76561198162340088"

var Session *discordgo.Session

var GetDB = make(chan *gorm.DB)
var DoneDB = make(chan bool)

func DBLoop(DB *gorm.DB) {
	for {
		GetDB <- DB
		<-DoneDB
	}
}

func RecoverPanic(channelID string) {

	if r := recover(); r != nil {

		s := Session

		// get the stack trace of the panic
		tempbuf := make([]byte, 10000)
		buflength := runtime.Stack(tempbuf, false)
		var buf []byte
		if buflength >= 1900 {
			buf = make([]byte, 1900)
		} else {
			buf = make([]byte, buflength)
		}
		runtime.Stack(buf, false)

		LogEvent(fmt.Sprintf("Recovering from panic: %v\n Stack trace: %s", r, buf), "ERR")
		if channelID != "" {
			s.ChannelMessageSend(channelID, "Error processing command.\nBug report sent to developers.")
		}

	}

}

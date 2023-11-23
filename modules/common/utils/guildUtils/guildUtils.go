package guildUtils

import (
	"github.com/angelfluffyookami/HSVRUSB/modules/common/global"
	"github.com/tjarratt/babble"
)

// Babbler babbles randomly generated keys for any sort of secure keys you need.
// wordCount int = amount of words to use for the key.
// Every word is 6 characters long,
func BabbleWords(wordCount int) string {
	var wordlist []string
	wordlist = append(wordlist, "grudge", "linear", "burial", "latest", "screen", "desert", "expose", "endure", "estate", "master", "refund", "throat", "effort", "pepper", "budget", "revive", "breast", "school", "flower", "ladder", "chorus", "wonder", "cheese", "sticky", "spread", "tumble", "vacuum", "flavor", "suntan", "mutter", "center", "punish", "resort", "hunter", "galaxy", "charge", "depend", "cotton", "shiver", "afford", "agenda", "timber", "morale", "behave", "camera", "expand", "carbon", "dollar", "latest", "mature", "mobile", "injury", "ensure", "barrel", "finish", "rhythm", "crutch", "museum", "lesson", "follow", "please", "safety", "modest", "remind", "reader", "demand", "ethics", "pledge", "accept", "ballot", "doctor", "gutter", "planet", "launch", "makeup", "freeze", "acquit", "colony", "rescue", "defend", "facade", "vision", "honest", "retire", "arrest", "banner", "thesis", "weight", "turkey", "worker", "column", "ignite", "facade", "ribbon", "bloody", "sacred", "inside", "dilute", "gallon", "theory", "behead", "proper", "chance", "single", "object", "temple", "modest", "likely", "adjust", "pastel", "attack", "market", "bishop", "belong", "effort", "rotate", "senior", "infect", "locate", "secure", "earwax", "normal", "flower", "prayer", "endure", "injury", "avenue", "family", "desert", "packet", "series", "tiptoe", "tumble", "harass", "spider", "output", "mutter", "church", "glance", "throne", "salmon", "option", "apathy", "cancer", "labour", "stroke", "dinner", "lounge", "gallon", "mobile", "bubble", "trance", "matrix", "ground", "escape", "defeat", "effect", "acquit", "square", "bitter", "excuse", "review", "normal", "formal", "player", "quaint", "belief", "critic", "accent", "empire", "junior", "lesson", "tongue", "voyage", "basket", "launch", "mosaic", "column", "margin", "source", "spirit", "cherry", "height", "bother", "deadly", "marble", "virtue", "devote", "mosque", "morale", "likely", "branch", "offend", "family", "script", "medium", "course", "theory", "weight", "winner")
	babbler := babble.NewBabbler()
	babbler.Count = wordCount
	babbler.Words = wordlist
	key := babbler.Babble()
	return key
}

// Gets the name of a guild.
// GID string guild ID of the server to target.
func GetGuildName(GID string) string {
	s := global.Session
	guild, err := s.Guild(GID)
	if err != nil {
		return "Undefined. " + GID
	} else {
		return guild.Name + " " + GID
	}
}

// Gets the username of the account that owns the specified server.
// GID string guild ID of the server to target.
func GetGuildOwnerName(GID string) string {
	s := global.Session
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

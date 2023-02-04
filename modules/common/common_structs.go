package common

import (
	"time"
)

type Data struct {
	Token         string
	AdminServer   string
	AdminChannel  string
	InfoChannel   string
	WarnChannel   string
	ErrChannel    string
	UpdateChannel string
	APIEndpoint   string
	Debugging     bool
}

type LogEntry struct {
	Time    time.Time
	Message string
	Level   string
}

var AuthKeyUpdater []func(GID string, AuthKey string, OldKey string)

type aircraftStruct struct {
	AV42c   int
	FA26b   int
	F45A    int
	AH94    int
	Invalid int
}

var Aircraft = aircraftStruct{
	AV42c:   1,
	FA26b:   2,
	F45A:    3,
	AH94:    4,
	Invalid: 0,
}

type weaponStruct struct {
	Gun     int
	AIM120  int
	AIM9    int
	AIM7    int
	AIM9X   int
	AIRST   int
	HARM    int
	AIM9E   int
	Invalid int
}

var Weapon = weaponStruct{
	Gun:     1,
	AIM120:  2,
	AIM9:    3,
	AIM7:    4,
	AIM9X:   5,
	AIRST:   6,
	HARM:    7,
	AIM9E:   8,
	Invalid: 0,
}

type teamStruct struct {
	Allied  int
	Enemy   int
	Invalid int
}

var Team = teamStruct{
	Allied:  1,
	Enemy:   2,
	Invalid: 0,
}

package bvr

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

type Online struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Team string `json:"team"`
}

type Kill struct {
	KillerID       string `json:"killerId,omitempty"`
	VictimID       string `json:"victimId,omitempty"`
	VictimTeam     int    `json:"victimTeam,omitempty"`
	KillerTeam     int    `json:"killerTeam,omitempty"`
	Time           int64  `json:"time,omitempty"`
	KillerAircraft int    `json:"killerAircraft,omitempty"`
	VictimAircraft int    `json:"victimAircraft,omitempty"`
	Weapon         int    `json:"weapon,omitempty"`
	ID             string `json:"id,omitempty"`
}

type Death struct {
	VictimID       string
	Time           int64
	VictimAircraft int
	ID             string
}

type LimitedUserData struct {
	ID         string   `json:"id"`
	PilotNames []string `gorm:"serializer:json" json:"pilotNames"`
	Kill       int      `json:"kills"`
	Deaths     int      `json:"deaths"`
	ELO        float64  `json:"elo"`
	Rank       int      `json:"rank"`
}

type User struct {
	UID         string
	ID          string             `json:"id"`
	PilotNames  []string           `gorm:"serializer:json" json:"pilotNames"`
	LoginTimes  []int64            `gorm:"serializer:json" json:"loginTimes"`
	LogoutTimes []int64            `gorm:"serializer:json" json:"logoutTimes"`
	Kills       int                `json:"kills"`
	Deaths      int                `json:"deaths"`
	Spawns      SpawnStruct        `gorm:"serializer:json" json:"spawns"`
	ELO         float64            `json:"elo"`
	ELOHistory  []ELOHistoryStruct `gorm:"serializer:json" json:"eloHistory"`
	Rank        int                `json:"rank"`
}

type ELOHistoryStruct struct {
	ELO  float64 `json:"elo"`
	Time int64   `json:"time"`
}

type SpawnStruct struct {
	AV42c   int `json:"0"`
	F26b    int `json:"1"`
	F45A    int `json:"2"`
	AH94    int `json:"3"`
	T55     int `json:"4"`
	Invalid int `json:"5"`
}

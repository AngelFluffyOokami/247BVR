package bvr

type AircraftStruct struct {
	AV42c   int
	FA26b   int
	F45A    int
	AH94    int
	Invalid int
}

var AircraftConst = AircraftStruct{
	AV42c:   1,
	FA26b:   2,
	F45A:    3,
	AH94:    4,
	Invalid: 0,
}

type WeaponStruct struct {
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

var WeaponConst = WeaponStruct{
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

type TeamStruct struct {
	Allied  int
	Enemy   int
	Invalid int
}

var TeamConst = TeamStruct{
	Allied:  1,
	Enemy:   2,
	Invalid: 0,
}

type OnlineStruct struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Team string `json:"team"`
}

type KillStruct struct {
	ID string `json:"id"`

	Killer       KillerStruct     `json:"killer"`
	VictimStruct VictimStruct     `json:"victim"`
	ServerInfo   ServerInfoStruct `json:"serverInfo"`

	Time int64 `json:"time"`
}

type ServerInfoStruct struct {
	MissionID   string   `json:"missionId"`
	OnlineUsers []string `json:"onlineUsers"`
	TimeOfDay   int      `json:"timeOfDay"`
}

type KillerStruct struct {
	OwnerID      string   `json:"ownerId"`
	Occupants    []string `json:"occupants"`
	AircraftType int      `json:"type"`
	Team         int      `json:"team"`
}

type VictimStruct struct {
	OwnerID      string   `json:"ownerId"`
	Occupants    []string `json:"occupants"`
	AircraftType int      `json:"type"`
	Team         int      `json:"team"`
}
type DeathStruct struct {
	ID string `json:"id"`

	Season     int              `json:"season"`
	ServerInfo ServerInfoStruct `json:"serverInfo"`
	Victim     VictimStruct     `json:"victim"`

	Time int64 `json:"Time"`
}

type LimitedUserDataStruct struct {
	ID         string   `json:"id"`
	PilotNames []string `json:"pilotNames"`
	Kill       int      `json:"kills"`
	Deaths     int      `json:"deaths"`
	ELO        float64  `json:"elo"`
	Rank       int      `json:"rank"`
}

type UserStruct struct {
	UID         string
	ID          string             `json:"id"`
	PilotNames  []string           `json:"pilotNames"`
	LoginTimes  []int64            `json:"loginTimes"`
	LogoutTimes []int64            `json:"logoutTimes"`
	Kills       int                `json:"kills"`
	Deaths      int                `json:"deaths"`
	TeamKills   int                `json:"teamKills"`
	Banned      bool               `json:"isBanned"`
	DiscordID   string             `json:"discordId"`
	Spawns      SpawnStruct        `json:"spawns"`
	ELO         float64            `json:"elo"`
	ELOHistory  []ELOHistoryStruct `json:"eloHistory"`
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

type EndOfSeasonStruct struct {
	Season    int     `json:"season"`
	ELO       float64 `json:"elo"`
	Rank      int     `json:"rank"`
	TeamKills int     `json:"teamKills"`
	History   string  `json:"history"`
}
